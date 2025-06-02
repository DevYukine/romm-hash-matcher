package main

import (
	"go.uber.org/zap"
	env2 "romm-hash-matcher/internal/env"
	"romm-hash-matcher/internal/hash"
	"romm-hash-matcher/internal/logging"
	"romm-hash-matcher/internal/model"
	"romm-hash-matcher/pkg/hasheous"
	"romm-hash-matcher/pkg/playmatch"
	"romm-hash-matcher/pkg/romm"
	"strconv"
	"strings"
	"sync"
)

var hasheousNotMapped = 0
var hasheousNotFound = 0
var hasheousIdentified = 0
var playmatchNotFound = 0
var playmatchFoundWithoutMetadata = 0
var playmatchIdentified = 0

func main() {
	logging.InitLogger()
	env2.InitEnv()

	rommClient := romm.NewClient(env2.Config.RommUrl, env2.Config.RommUsername, env2.Config.RommPassword)
	playmatchClient := playmatch.NewClient()
	hasheosClient := hasheous.NewClient()

	logging.Logger.Info("Fetching unmatched ROMs from RomM, this may take a while depending on the number of ROMs in your collection")

	roms, err := rommClient.GetUnmatchedRoms()
	if err != nil {
		logging.Logger.Fatal("Error fetching unmatched ROMs", zap.Error(err))
	}

	logging.Logger.Info("Found unmatched ROMs", zap.Int("count", len(*roms)))

	zipped := 0

	for _, rom := range *roms {
		if strings.HasSuffix(strings.ToLower(rom.Name), ".zip") {
			zipped++
		}
	}

	logging.Logger.Info("Found unmatched ROMs, searching by hash now", zap.Int("count", len(*roms)), zap.Int("zipped", zipped))

	var wg sync.WaitGroup

	currentCount := 0

	for _, rom := range *roms {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var metadataOfRom *model.InternalRom

			if !strings.HasSuffix(strings.ToLower(rom.Name), ".zip") {
				logging.Logger.Debug("Found not zipped Rom, use hash from RomM instead hashing it", zap.String("romName", rom.Name))

				if rom.Md5Hash == nil || rom.Sha1Hash == nil {
					logging.Logger.Warn("ROM does not have MD5 or SHA1 hash, please run the hashing scan before using this tool once", zap.String("romName", rom.Name))
					return
				}

				metadataOfRom = &model.InternalRom{
					Name: rom.Name,
					Size: int64(rom.FsSizeBytes),
					MD5:  *rom.Md5Hash,
					SHA1: *rom.Sha1Hash,
				}
			} else {
				logging.Logger.Debug("Found zipped Rom, downloading and hashing the content", zap.String("romName", rom.Name))
				if rom.FsSizeBytes > hash.MaxZipExtractionFileSize {
					logging.Logger.Warn("Found zipped Rom with size larger than 2GB, skipping for now as hashing would take a long time", zap.String("romName", rom.Name))
					return
				}

				metadataOfRom, err = rommClient.GetMetadataOfZippedRom("/api/roms/" + strconv.FormatInt(int64(rom.Id), 10) + "/content/" + rom.Name)
				if err != nil {
					logging.Logger.Error("Error getting metadata of zipped ROM", zap.Error(err), zap.String("romName", rom.Name))
					return
				}

				if metadataOfRom == nil {
					return
				}
			}

			logging.Logger.Debug("ROM metadata", zap.String("name", metadataOfRom.Name), zap.Int64("size", metadataOfRom.Size), zap.String("md5", metadataOfRom.MD5), zap.String("sha1", metadataOfRom.SHA1), zap.String("sha256", metadataOfRom.SHA256))

			igdbId, err := matchPlaymatch(playmatchClient, metadataOfRom)

			if err != nil {
				logging.Logger.Error("Error identifying game via Playmatch", zap.Error(err), zap.String("romName", rom.Name))
			}

			slug, err := matchHasheous(hasheosClient, metadataOfRom)

			if err != nil {
				logging.Logger.Error("Error looking up ROM by hash on Hasheous", zap.Error(err), zap.String("romName", rom.Name))
			}

			var searchResult *romm.SearchResponse

			isPlaymatchMatch := false

			if igdbId != nil {
				logging.Logger.Debug("Searching for igdb result by Playmatch Id", zap.Int64("igdbId", *igdbId))
				searchResult, err = rommClient.SearchMetadataByIgdbId(rom.Id, *igdbId)

				if err != nil {
					logging.Logger.Error("Error searching metadata by IGDB ID", zap.Error(err), zap.String("romName", rom.Name), zap.Int64("igdbId", *igdbId))
					return
				}

				isPlaymatchMatch = true
			} else if slug != nil {
				logging.Logger.Debug("Searching for igdb result by Hasheous Slug", zap.String("slug", *slug))

				igdbId, err = playmatchClient.GetIdBySlug(*slug)

				if err != nil {
					logging.Logger.Error("Error getting IGDB ID from playmatch by Hasheous slug", zap.Error(err), zap.String("romName", rom.Name), zap.String("slug", *slug))
					return
				}

				searchResult, err = rommClient.SearchMetadataByIgdbId(rom.Id, *igdbId)

				if err != nil {
					logging.Logger.Error("Error searching metadata by IGDB ID", zap.Error(err), zap.String("romName", rom.Name), zap.Int64("igdbId", *igdbId))
					return
				}
			} else {
				logging.Logger.Info("Could not match rom to IGDB id by hash, consider contributing to Playmatch or Hasheous and match it for the community!", zap.String("romName", rom.Name))
				return
			}

			if searchResult == nil || searchResult.IgdbId == nil {
				logging.Logger.Warn("No IGDB ID found in search result", zap.String("romName", rom.Name), zap.Int64p("igdbId", igdbId), zap.Stringp("slug", slug))
				return
			}

			err = rommClient.ManuallyMatchRom(rom.Id, rom.FsName, searchResult)

			if err != nil {
				logging.Logger.Error("Error manually matching ROM", zap.Error(err), zap.String("romName", rom.Name), zap.Int64("romId", rom.Id), zap.Int64p("igdbId", igdbId))
				return
			}

			logging.Logger.Info("Successfully matched ROM in RomM via Hash", zap.String("romName", rom.Name), zap.Bool("IsMatchedViaPlaymatch", isPlaymatchMatch), zap.Bool("IsMatchedViaHasheous", !isPlaymatchMatch), zap.Int64("romId", rom.Id), zap.Int64p("igdbId", igdbId), zap.Stringp("igdbSlug", slug), zap.String("igdbName", searchResult.Name), zap.String("igdbUrlCover", searchResult.IgdbUrlCover))
		}()

		currentCount++
	}

	wg.Wait()
	logging.Logger.Info("Finished processing unmatched ROMs", zap.Int("playmatchNotFound", playmatchNotFound), zap.Int("playmatchFoundWithoutMetadata", playmatchFoundWithoutMetadata), zap.Int("playmatchIdentified", playmatchIdentified), zap.Int("hasheousNotMapped", hasheousNotMapped), zap.Int("hasheousNotFound", hasheousNotFound), zap.Int("hasheousIdentified", hasheousIdentified))
}

func matchPlaymatch(client *playmatch.Client, rom *model.InternalRom) (*int64, error) {
	identifyResponse, err := client.IdentifyGame(rom)

	if err != nil {
		logging.Logger.Error("Error identifying game", zap.Error(err), zap.String("romName", rom.Name))
		return nil, err
	}

	if identifyResponse == nil {
		logging.Logger.Error("Playmatch returned nil response", zap.String("romName", rom.Name))
		return nil, nil
	}

	if identifyResponse.GameMatchType == playmatch.GameMatchTypeNoMatch || identifyResponse.ExternalMetadata == nil || len(identifyResponse.ExternalMetadata) == 0 {
		logging.Logger.Debug("Playmatch no match found for game", zap.String("romName", rom.Name))
		playmatchNotFound++
		return nil, nil
	}

	igdbMetadata := identifyResponse.ExternalMetadata[0]

	if igdbMetadata.MatchType == playmatch.MetadataMatchTypeFailed || igdbMetadata.MatchType == playmatch.MetadataMatchTypeNone {
		logging.Logger.Debug("Playmatch match found but no igdb match is present", zap.String("romName", rom.Name), zap.String("matchType", string(igdbMetadata.MatchType)), zap.String("failedMatchReason", string(igdbMetadata.FailedMatchReason)))
		playmatchFoundWithoutMetadata++
		return nil, nil
	}

	playmatchIdentified++
	logging.Logger.Debug("Playmatch Identified Game", zap.String("romName", rom.Name), zap.String("id", identifyResponse.Id), zap.String("igdbId", igdbMetadata.ProviderId), zap.String("matchType", string(igdbMetadata.MatchType)), zap.String("automaticMatchReason", string(igdbMetadata.AutomaticMatchReason)), zap.String("comment", igdbMetadata.Comment))
	id, err := strconv.ParseInt(igdbMetadata.ProviderId, 10, 64)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

func matchHasheous(client *hasheous.Client, rom *model.InternalRom) (*string, error) {
	lookupResponse, err := client.LookupByHash(*rom)

	if err != nil {
		logging.Logger.Error("Error looking up ROM by hash on Hasheous", zap.Error(err), zap.String("romName", rom.Name))
		return nil, err
	}

	if lookupResponse == nil {
		logging.Logger.Debug("No match found on Hasheous", zap.String("romName", rom.Name))
		return nil, nil
	}

	var igdbMedata *hasheous.Metadata

	for _, metadata := range lookupResponse.Metadata {
		if strings.ToLower(metadata.Source) == "igdb" {
			igdbMedata = &metadata
		}
	}

	if igdbMedata == nil {
		logging.Logger.Debug("no IGDB metadata found on Hasheous", zap.String("romName", rom.Name))
		hasheousNotFound++
		return nil, nil
	}

	if igdbMedata.Status == "NotMapped" {
		logging.Logger.Debug("Hasheous IGDB metadata found but not mapped", zap.String("romName", rom.Name), zap.String("status", igdbMedata.Status))
		hasheousNotMapped++
		return nil, nil
	}

	hasheousIdentified++
	logging.Logger.Debug("Hasheos identified Game", zap.String("romName", rom.Name), zap.String("igdbSlug", igdbMedata.Id), zap.String("status", igdbMedata.Status), zap.String("matchType", igdbMedata.MatchMethod))

	return &igdbMedata.Id, nil
}

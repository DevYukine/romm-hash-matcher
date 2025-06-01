# RomM Hash Matcher

This is a small program i threw together to match hashes of roms in your RomM Server to a database of known hashes, RomM developer plan on adding this as a native feature but until then you can use this and manually do it once to match unmatched RomM automatically (this works best when your rom is matching an entry from known DAT Groups like No-Intro, Redump, TOSEC etc)

### Only tested with IGDB so far

## What does it do?

1. Get list of unmatched roms from your RomM Server
2. Get hashes of those roms (if RomM has the hashes already it will use those, for .zip files it will download and unzip it and then get the hashes of the file inside as DAT groups do not hash .zip files but only the roms inside)
3. Do a lookup of those hashes on [Playmatch](https://github.com/RetroRealm/playmatch) and [Hasheous](https://hasheous.org/) via their API's
4. If a match on either is found, it will search for metadata on IGDB and update the RomM Server with the metadata (this does the same as when you do a manual search in the UI)

## Docs

### Environment Variables

You need to set following environment variables for the program to work:

- `ROMM_USERNAME` - Your RomM username
- `ROMM_PASSWORD` - Your RomM password
- `ROMM_URL` - The URL of your RomM server (e.g. `http://localhost:8000`) including the protocol!

You can also use a .env file to set these variables, just create a `.env` file in the same directory as the script and add the variables there.

# TenPack - Velocidrone Laptimer and Telemetry (Local) 

## A Velocidrone companion app that allows for personal tracking and the tracking of other pilot's laps + telemetry.

### Status: Working / Incomplete
This is a fun project I took on as I continue to learn as a hobbyist. The app is a lap time that also provides telemetry. The laptimes are in order leading with the fastest final times at top. Green in any lap column indicates who had the fastest time of that lap. Splits are based off the fastest overall time in the reported race which is displayed in blue. Split times will show as red for being slower than than the respective split section of the fastest overall lap, or green for faster.  The name of the track and the desired split sections can be specified in the settings JSON. If the highest value gate ID exceeds the track gate length the telemtry will default to showing splits between every gate. The websocket binds to your local ip I.e 192.168.68.55 style ip and not 127.0.0.1. This is an unfinished project and not meant to be taken serious at all. Enjoy

Things that are broken:
- Abandoned races will not register as non-race host since the VD API currently doesn't send messages for it
- Specifying gates out of the attempted tracks range (**fixed** - will improve handling later)
- other stuff

![Image of Laptimes with split times](https://raw.githubusercontent.com/k1tig/TenPack/refs/heads/main/WS/TenpackCLI.png)

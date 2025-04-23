**NOTE:** If you are upgrading from v0.1.0, which is the first version, please refer to [upgrade-from-v0.1.0.md](./upgrade-from-v0.1.0.md)

```bash
# Ensure compose yml up to date.
## Mac/Linux
curl -L -O https://raw.githubusercontent.com/glidea/zenfeed/main/docker-compose.yml
## Windows PowerShell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/glidea/zenfeed/main/docker-compose.yml" -OutFile ([System.IO.Path]::GetFileName("https://raw.githubusercontent.com/glidea/zenfeed/main/docker-compose.yml"))


# Ensure images up to date.
docker-compose -p zenfeed pull


# Upgrading without reconfiguring, etc APIKey.
docker-compose -p zenfeed up -d
```

Then all the feed data and configurations should be intact.

# Fake-AI-Detective

[![Docker Image CI](https://github.com/wzxjohn/fake-ai-detective/actions/workflows/docker-image.yml/badge.svg)](https://github.com/wzxjohn/fake-ai-detective/actions/workflows/docker-image.yml)

Fake AI Detective is a tool to test your OpenAI like API channel.

# Usage

## Deploy

Cloudflare CDN is recommended.

### Docker Composer
```
wget https://raw.githubusercontent.com/wzxjohn/fake-ai-detective/refs/heads/master/docker-compose.yaml
# Change DETECTIVE_DOMAIN env in config
docker-compose up -d
```

## Config environments

### `DETECTIVE_DOMAIN`

Define the domain of the tool. Must be accessible through the Internet.

### `DETECTIVE_API_PREFIX`

Define the API Prefix of the tool. Set if you want to deploy under sub folder.

### `DETECTIVE_IMAGE_PREFIX`

Define the Image Prefix of the tool. Set if you want to deploy under sub folder.

## Notice

If you set the `DETECTIVE_API_PREFIX` or `DETECTIVE_IMAGE_PREFIX` env, please make sure your reserve proxy has been
properly configured.

# Stargazers over time

[![Stargazers over time](https://starchart.cc/wzxjohn/fake-ai-detective.svg)](https://starchart.cc/wzxjohn/fake-ai-detective)
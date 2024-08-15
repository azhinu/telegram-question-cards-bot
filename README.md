# Telegram question cards bot

This bot will send you a question card from a list of questions that you can load from a YAML file.

## Usage

Bot can run in two modes: polling and webhook. Polling is the default mode and is used when no webhook URL is provided. Default webhook port is 1443.

### Cli flags:
  -h      | --help            | Env vars            | Show context-sensitive help.       
----------|-------------------|---------------------|----------------------------------
  &nbsp;  | --version         | &nbsp;              | Print version information and quit 
  -d      | --debug           | QC_BOT_DEBUG        | Enable debug logging            
  -t      | --token=STRING    | QC_BOT_TOKEN        | Telegram bot token. Required    
  -u      | --url=https://example.com/bot-secret-url   | QC_BOT_URL       | Webhook URL. Optional    
  -p      | --port=1443       | QC_BOT_PORT         | Weebhook port

## Question Decks

You can load a list of questions from a YAML file. The file should be structured as follows:

```yaml
deck1:
  - "q1"
  - "q2"
  - "q3"
  - "q4"
deck2:
  - "q1"
  - "q2"
```

# Deployment Tool
This tool starts the server on the port of your choice and waits for a Webhook from GitHub.
When it receives the webhook, the tool checks the validity of the payload as well as the signature.
Once the checks are successful, the tool starts an action thread.

Each action can have its own condition, depending on which action will be executed or not
(the output codes are controlled: 0 - command is executed, otherwise - not).

The tool guarantees the execution of actions in the order in which they are listed in the configuration file.
The tool guarantees that webhooks are processed in the order in which they are received (one after the other).

## Smart payload access
You can access webhook payloads with [GJSON (link to documentation)](docs/GJSON.md). GJSON syntax is available in the `if` and `run` fields
(use curly braces to process gjson). For example `repo name: {repository.name}`
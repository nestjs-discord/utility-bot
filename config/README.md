# Config

The `config` package is essential for setting up and maintaining the application's configuration.

It ensures that the application is configured correctly and ready to run.
As the `README.md` file in the root directory outlines, this project relies on environment variables
and a YAML file to configure its behavior.

This package reads `.env` and `config.yaml` files from the disk, parses them, and validates them.
Any error thrown from this package should halt the application's bootstrap process.

In addition to managing these files, it also holds the configuration related to the Discord bot, including permissions and intents.

# LangForge: A Toolkit for Creating and Deploying LangChain Apps

LangForge is an open-source toolkit designed to make it easy to create and deploy LangChain applications.
With LangForge, you can easily set up your environment, manage API keys, chat with your chains using the Jupyter integration, and automatically generate a REST interface for your app. Additionally, LangForge comes with a collection of predefined notebooks tailored for various use cases, allowing you to quickly get started and adapt them to your specific needs.

## Features

- Simplified environment setup and API key management
- Predefined notebooks for various use cases to help you get started quickly
- Instantly chat with your chains using the Jupyter integration
- Automatic REST interface generation for your app

## Installation

To install LangForge, simply run the following command:

```bash
pip install langforge-cli
```

## Getting Started

1. **Create a new LangChain app:** Use the create command to generate a new LangChain app. LangForge will set up a virtual environment, install the required packages, and configure the API keys, settingup a ready-to-use foundation for building your app.

```bash
langforge create myapp
```

2. **Launch JupyterLab:** Run the langforge lab command to launch Jupyter Lab. LangForge Lab comes with ready-to-use notebooks for various use cases and an integration that allows you to chat with your chains directly within Jupyter.

```bash
langforge lab
```

3. **Serve your app:** LangForge automatically generates a REST interface for your app, making it easy to deploy and share with others. Use the `serve` command followed by the name of your notebook to start serving your app.

```bash
langforge serve my_notebook.ipynb
```

## Contributing

We welcome contributions from the community! If you'd like to contribute to LangForge, please feel free to submit pull requests or open issues on our GitHub repository.

## License

LangForge is released under the MIT License.

# Spacy

This module uses spacy to compute similarities between tasks.

## Installation

*It is developped with python 3.6*

```bash
virtualenv env
source env/bin/activate
pip install -U pip # spacy needs an updated pip to download its models
pip install flask spacy
python -m spacy download en
```

## Starting the server

```bash 
python main.py
```

This will start the server on port 1717. Update `main.py` to:
- disable the debug mode
- load another spaCy model
- do other things :stuck_out_tongue:

## Trouble shooting

The error about `--no-cache-dir` is solved by running the `pip install -U pip` command.

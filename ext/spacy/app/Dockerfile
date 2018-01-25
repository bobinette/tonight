FROM python:3.6

# Ensure the last pip is installed
RUN pip install -U pip
RUN pip install flask spacy

# Use the basic en model. If you change it here, make
# sure you also change it in main.py
RUN python -m spacy download en

# Copy app
COPY main.py main.py

ENTRYPOINT ["python"]
CMD ["main.py"]

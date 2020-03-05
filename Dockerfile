FROM python:3.7
RUN mkdir /code
COPY ./requirements.txt /code/
RUN pip install -r /code/requirements.txt
COPY . /code/
WORKDIR /code

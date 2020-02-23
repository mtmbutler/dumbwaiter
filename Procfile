release: python manage.py migrate --noinput && python manage.py createsuperuser --username admin --email "" --noinput
web: gunicorn simplefi.wsgi --log-file -

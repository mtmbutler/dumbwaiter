release: python manage.py migrate --noinput && python manage.py createsuperuser --username admin --email admin@admin.com --noinput
web: gunicorn server.wsgi --log-file -

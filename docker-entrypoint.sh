python manage.py collectstatic --no-input
python manage.py migrate
gunicorn server.wsgi -b 0.0.0.0:8000

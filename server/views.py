from django.contrib.auth.models import User
from rest_framework import viewsets

from server.serializers import UserSerializer


class UserViewSet(viewsets.ModelViewSet):
    serializer_class = UserSerializer

    def get_queryset(self):
        user = self.request.user
        if user.is_staff or user.is_superuser:
            return User.objects.all()
        elif user.is_authenticated:
            return User.objects.filter(id=user.id)
        return User.objects.none()

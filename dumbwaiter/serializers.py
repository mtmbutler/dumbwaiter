from rest_framework import serializers

from dumbwaiter import models as m


class DaySerializer(serializers.HyperlinkedModelSerializer):
    user = serializers.PrimaryKeyRelatedField(
        read_only=True, default=serializers.CurrentUserDefault()
    )

    class Meta:
        model = m.Day
        fields = [
            "user",
            "date",
            "am_weight",
            "pm_weight",
            "snack",
            "breakfast",
            "lunch",
            "dinner",
            "exercise",
        ]
        read_only_fields = ["user"]

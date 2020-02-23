from rest_framework import serializers

from dumbwaiter import models as m


class DaySerializer(serializers.HyperlinkedModelSerializer):

    class Meta:
        model = m.Day
        fields = [
            "date",
            "am_weight",
            "pm_weight",
            "snack",
            "breakfast",
            "lunch",
            "dinner",
            "exercise",
        ]



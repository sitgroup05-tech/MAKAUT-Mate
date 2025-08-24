
from django.urls import path

from . import views

urlpatterns = [
    path('',view= views.home , name='home'),
    path('hello',views.home,name='hello'),
]
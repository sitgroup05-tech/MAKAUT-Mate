
from django.contrib import admin
from django.urls import path , include

from . import views


from django.conf import settings
from django.conf.urls.static import static

urlpatterns = [
    path('admin/', admin.site.urls),
    # path('', views.home , name='home'),
    path('',include('home.urls')),
    path('home',include('home.urls'),name = 'home' ),


    path('about',views.home,name = 'about' ),
    path('courses',views.home,name = 'courses' ),
    path('contact',views.home,name = 'contact' ),
    path('privacy',views.home,name = 'privacy' ),
    path('terms',views.home,name = 'terms' ),
]


if settings.DEBUG:
    urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
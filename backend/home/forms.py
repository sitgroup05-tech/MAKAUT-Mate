from django import forms
from .models import UploadedFile 

from django.db import models

class UploadedFileForm(forms.ModelForm):
    class Meta:
        model = UploadedFile
        fields = ['name', 'email', 'file']


class SearchDocument(forms.Form):

    context = forms.CharField(
        label='Seach :',
        widget=forms.TextInput(
            attrs={
                'class' : 'mx-1 my-auto',
                'placeholder' : 'Search Document'
            }
        )

    )


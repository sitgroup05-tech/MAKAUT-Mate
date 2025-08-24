from django.db import models

class UploadedFile(models.Model):
    name = models.CharField(max_length=100)
    email = models.EmailField()
    file = models.FileField(upload_to='uploads/')  # will be saved in MEDIA_ROOT/uploads
    uploaded_at = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return f"{self.name} - {self.file.name}"
    


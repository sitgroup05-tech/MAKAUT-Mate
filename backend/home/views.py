from django.shortcuts import render, redirect
from .models import UploadedFile 

import background
import pathlib

from TextProcessor import tokenize

from . forms import SearchDocument
from .forms import UploadedFileForm


import easyocr

reader = easyocr.Reader(['en'])


class Result :
    have_result : bool = True


# note this can handel low powered background task
@background.task
def process_file(file_path):
    """
        process the text from the file
    """
    print(file_path)
    path = pathlib.Path(file_path)
    if path.suffix in ['.jpeg','.png','.jpg'] : 
        results = reader.readtext(file_path)
        extracted_text = " ".join([res[1] for res in results])
    
    elif path.suffix == '.txt' :
      with path.open() as fp : 
          extracted_text = fp.read()
    print('text : ' , extracted_text)
    vector_data = tokenize.process_text(extracted_text)

    print('SIZE : ' ,vector_data.shape)

def home(request):
    files = UploadedFile.objects.all().order_by('-uploaded_at')
    result = Result()
    search = SearchDocument()

    action = request.GET.get('action')

    if request.method == 'POST':
        form = UploadedFileForm(request.POST, request.FILES)
        if form.is_valid():
            # Perform OCR
            doc = form.save()
            process_file(doc.file.path)
            
            return redirect('home')
    else: # GET
        form = UploadedFileForm()

    return render(request, 'pages/home.html', {'files': files, 'form': form , 'search' : search , 'result' : result })


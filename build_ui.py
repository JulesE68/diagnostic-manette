import base64,io
from PIL import Image
s=open('ui/src.html').read()
css=open('../deps/7css-pages/7.css').read()
im=Image.open('ui/icon.png').resize((48,48),Image.LANCZOS);b=io.BytesIO();im.save(b,'PNG')
s=s.replace('/*7CSS*/',css,1).replace('/*ICON*/','data:image/png;base64,'+base64.b64encode(b.getvalue()).decode(),1)
open('ui/app.html','w').write(s)
print(len(s))

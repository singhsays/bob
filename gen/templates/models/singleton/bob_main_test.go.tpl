{{$doneTypes := dict }}
{{- range $table := .Tables}}
{{- $tAlias := $.Aliases.Table $table.Key}}
  {{range $column := $table.Columns -}}
    {{- $colTyp := $column.Type -}}
    {{- if hasKey $doneTypes $colTyp}}{{continue}}{{end -}}
    {{- $_ :=  set $doneTypes $colTyp nil -}}
    {{- $typInfo :=  index $.Types $column.Type -}}
    {{- if $typInfo.NoScannerValuerTest}}{{continue}}{{end -}}
    {{- if isPrimitiveType $colTyp}}{{continue}}{{end -}}
    {{- $.Importer.ImportList $typInfo.Imports -}}
      {{$.Importer.Import "database/sql"}}
      {{$.Importer.Import "database/sql/driver"}}
      // Make sure the type {{$colTyp}} satisfies database/sql.Scanner
      var _ sql.Scanner = (*{{$colTyp}})(nil)

      // Make sure the type {{$colTyp}} satisfies database/sql/driver.Valuer
      var _ driver.Valuer = *new({{$colTyp}})

  {{end -}}
{{- end}}

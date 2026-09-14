generator target {
  output_type "rest/openapi"
  output_dir "../"
  version "0.226.2"

  on_generate "openapi-generator-cli generate -i openapi.yaml -g go -o . --additional-properties enumClassPrefix=true,isGoSubmodule=true,packageName=baml_client,withGoMod=false,generateModelTests=false,generateApiTests=false && rm -rf test"
}
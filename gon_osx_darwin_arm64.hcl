source = ["./dist/astra-cli-osx-arm_darwin_arm64/astra-cli"]
bundle_id = "pro.foundev.astra.cli"

sign {
  application_identity = "DDF52A8D387B7E77F18F86043FF2AC7AED277179"
}

zip {
  output_path = "./dist/astra-cli-osx-arm-signed.zip"
}

apple_id {
  password = "@env:AC_PASSWORD"
}

source = ["./dist/astra-osx-arm_darwin_arm64/astra"]
bundle_id = "pro.foundev.astra.cli"

sign {
  application_identity = "DDF52A8D387B7E77F18F86043FF2AC7AED277179"
}

zip {
  output_path = "./dist/astra-osx-arm-signed.zip"
}

apple_id {
  password = "@env:AC_PASSWORD"
}

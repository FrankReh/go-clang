# Set the CGO_LDFLAGS based on the llvm-config location.

echo "before CGO_LDFLAGS       $CGO_LDFLAGS"
export CGO_LDFLAGS="-L`llvm-config --libdir`" 
echo "after  CGO_LDFLAGS       $CGO_LDFLAGS"

case $(uname) in
    Darwin)
        export SDKROOT="$(xcrun --sdk macosx --show-sdk-path)"

        echo "before DYLD_LIBRARY_PATH   $DYLD_LIBRARY_PATH"
        export DYLD_LIBRARY_PATH=$(llvm-config --libdir)
        echo "after  DYLD_LIBRARY_PATH   $DYLD_LIBRARY_PATH"

        echo "before CGO_CPPFLAGS   $CGO_CPPFLAGS"
        CGO_CPPFLAGS="-I /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include"
        CGO_CPPFLAGS+=" -Wno-expansion-to-defined"
        CGO_CPPFLAGS+=" -Wno-nullability-completeness"
        CGO_CPPFLAGS+=" -Wno-undef-prefix"
        export CGO_CPPFLAGS

        echo "after  CGO_CPPFLAGS   $CGO_CPPFLAGS"

        alias cc=clang
        ;;
    Linux|FreeBSD)
        echo "before LD_LIBRARY_PATH   $LD_LIBRARY_PATH"
        export LD_LIBRARY_PATH=$(llvm-config --libdir)
        echo "after  LD_LIBRARY_PATH   $LD_LIBRARY_PATH"
        ;;
esac

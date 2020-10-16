# Set the CGO_LDFLAGS based on the llvm-config location.

echo "before CGO_LDFLAGS       $CGO_LDFLAGS"
export CGO_LDFLAGS="-L`llvm-config --libdir`" 
echo "after  CGO_LDFLAGS       $CGO_LDFLAGS"

case $(uname) in
    Darwin)
        echo "before DYLD_LIBRARY_PATH   $DYLD_LIBRARY_PATH"
        export DYLD_LIBRARY_PATH=$(llvm-config --libdir)
        echo "after  DYLD_LIBRARY_PATH   $DYLD_LIBRARY_PATH"

        echo "before CGO_CPPFLAGS   $CGO_CPPFLAGS"
        export CGO_CPPFLAGS="-I /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include -Wno-expansion-to-defined -Wno-nullability-completeness"
        echo "after  CGO_CPPFLAGS   $CGO_CPPFLAGS"

        alias cc=clang
        ;;
    Linux|FreeBSD)
        echo "before LD_LIBRARY_PATH   $LD_LIBRARY_PATH"
        export LD_LIBRARY_PATH=$(llvm-config --libdir)
        echo "after  LD_LIBRARY_PATH   $LD_LIBRARY_PATH"
        ;;
esac

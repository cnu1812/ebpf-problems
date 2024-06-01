## Install the required dependencies:

``` go get github.com/cilium/ebpf ```

## Compiling the eBPF code

``` go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cgo drop_tcp_port.c -- -I$DIR/../headers ```

or 

``` go generate .```

## Compiling the go code

``` go build```

## Runnig the code

``` ./drop_tcp_port [Enter the port numner of your choice] ```


Assuming all the pre-required dependencies are installed. 

Here I'm not uploading the auto generated files, to get them enter this command ``` go generate .```

Reference: https://youtu.be/Hed2DOrk_kk?si=XgVK3fcYOcgG4-Rf


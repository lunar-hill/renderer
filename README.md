# FauxGL OCI FaaS Renderer 

we replaced the entire renderer with brick hill's, so more to do :sob:

## how to test

1. install dockah (docker)

2. install fn
```sh
curl -LSs https://raw.githubusercontent.com/fnproject/cli/master/install | sh
```

3. run fn server (on a seperate terminal or window idk)
```sh
fn start
```

4. create app
```sh
fn create app goapp
```

5. deploy app
```sh
fn --verbose deploy --app goapp --local
```

6. test
```sh
echo -n '{"avatarJSON":"","size":512}' | fn invoke goapp render --content-type application/json
```

## Todo:
- [x] Re-add UUID lmfao 
- [x] Replace baker's rigs with brick hill rigs
- [x] Add t-shirt loading
- [ ] Add pre-cached resources (requires custom dockerfile)

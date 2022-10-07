# Running the example

Requirements:

- Make sure you have Dapr and the Dapr CLI installed. See Dapr's [Get started](https://docs.dapr.io/getting-started/) docs. You will need Dapr 1.9 or higher.
- You also need Go 1.19 or higher.

First, start the Dapr component. In the root folder of this project (one level up from here), run:

```sh
go run .
```

In **another terminal**, from this `example` folder, run:

```sh
# cd example

dapr run \
  --app-id myapp \
  --components-path ./components \
  --log-level debug \
  --\
      go run main.go
```

using Go = import "go.capnp";
@0x85d3acc39d94e0f8;

$Go.package("ping");
$Go.import("ping");

struct Ping {
  message @0 :Text;
}

struct BlockInfo{
    message @0 :Data;
}
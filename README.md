# gomap

## description

Gomap is golang implementation of [cgimap](https://github.com/zerebubuth/openstreetmap-cgimap).

## demo

Gomap demo is located [here](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/).

Implemented API call list:

* changesets:

  * GET /api/0.6/changeset/#id?include_discussion=true
    * [changeset 58719365](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/changeset/58719365)

* elements:

  * GET /api/0.6/[node|way|relation]/#id
    * [node 21140736](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/node/21140736)
    * [way 19780617](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/way/19780617)
    * [relation 16239](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/relation/16239)
  * GET /api/0.6/[node|way|relation]/#id/history
    * [node 21140736 history](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/node/21140736/history)
    * [way 19780617 history](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/way/19780617/history)
    * [relation 16239 history](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/relation/16239/history)
  * GET /api/0.6/[node|way|relation]/#id/#version
    * [node 21140736 version 10](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/node/21140736/10)
    * [way 19780617 version 56](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/way/19780617/56)
    * [relation 16239 version 1061](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/relation/16239/1061)
  * GET /api/0.6/[nodes|ways|relations]?#parameters
    * [node 21140736 and 21140802v3](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/nodes?nodes=21140736,21140802v3)
    * [way 19780617 and 24530399v7](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/ways?ways=19780617,24530399v7)
    * [relation 22868 and 27939v8](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/relations?relations=22868,27939v8)
  * GET /api/0.6/node/#id/ways
    * [ways for node 21140736](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/node/21140736/ways)
  * GET /api/0.6/[way|relation]/#id/full
    * [way 19780617 full](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/way/19780617/full)
    * [relation 16239 full](https://zkmeyj45t6.execute-api.us-west-2.amazonaws.com/staging/api/0.6/relation/16239/full)
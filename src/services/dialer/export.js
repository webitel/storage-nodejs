/**
 * Created by I. Navrotskyj on 21.11.17.
 */

"use strict";

console.log("FORK");
throw 1;
process.on('message', (m) => {
    console.log('CHILD got message:', m);
});
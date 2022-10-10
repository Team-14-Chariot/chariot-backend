import Pocketbase from 'pocketbase';


/* Basic javascript 
- use var. const, or let to declare a variable, not types
ie. var x = 6

- strings can be single or double quotes

- function [function name] ([args]) {}
ie. function calcCelsius(fahrenheit) {
    return ...;
}

- onClick=[function name] to respond to a certain user input

- arrays are the same as Java

- JSON is like dictionaries in Python
ie. look below at how data is formatted
*/

// create new pocketbase server
const client = new Pocketbase('http://127.0.01.8090')

// create record entry to insert into organizers table
data = {username: "testuser@gmail.com", password: ""}

// submit request to insert data into table
// check if it worked by looking at admin UI dashboard
const record = await client.records.create('organizers', data)
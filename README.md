# POC BookMetaData
## Problem
When reading LitRPG, I often found myself wondering what the current stats 
blocks of the characters looked like. The more I thought about the feature 
the more I realized I liked the problem.
## Objective
Create a solution that takes in the readers current page and providers 
metadata to the reader. I wanted the solution to be generalized so that 
the solution would work for any metadata not just metadata for LitRPG.
## Notes
cinnamonBun was the first pass. The JSON structure is more rudimentary than
threadBare. When I went to create a real example, I realized that I would
need every key to have a display value. This is turn created display/value
structure.
```
"display": ""
"value": 
```
## Solution
There are two assumptions made.
1. No arrays are used in the JSON. This is to prevent two problems.
   1. Arrays in json are just maps with default keys that are the index of the elements. This creates a problem when those elements move around.
   2. There is no way to update elements in an array which creates much worse data entry.
2. `keys` are present in the root json. Current implementation has static set keys, but I believe this can be modified to by dynamic in future interations.  

In a full implementation a template system would also be necessary for rendering the output.
To enable this functionality to be integrated into other software it will need to be wrapped in a service.
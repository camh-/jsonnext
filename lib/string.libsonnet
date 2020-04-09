{
  local string = self,

  // containsChar returns true if @char is in @str, false if it is not. @char
  // must be a string of length one. Note: The implementation uses
  // `std.strReplace` as this is built-in in the Go implementation and is
  // likely the fastest way.
  containsChar(str, char)::
    if !std.isString(char) then
      error ('string.containsChar second param must be a string, got ' + std.type(char))
    else if std.length(char) != 1 then
      error ('string.containsChar second param must be length 1, got ' + std.length(char))
    else
      std.strReplace(str, char, '') != str,
}

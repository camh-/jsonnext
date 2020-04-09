{
  local array = self,
  local op = import 'op.libsonnet',
  local values = import 'values.libsonnet',

  // make returns its argument if it is an array, otherwise returns a single
  // element array containing the argument.
  make(any):: if std.isArray(any) then any else [any],

  // coalesce returns the first non-null value from the given array or null
  // if there are none.
  coalesce(arr)::
    if arr == [] then null
    else if arr[0] != null then arr[0]
    else array.coalesce(arr[1:]),

  // collapse returns the input array with any null values removed
  collapse(arr):: [e for e in arr if e != null],

  // accumulate adds all the elements in the array @arr together starting
  // with the value @init. If @init is not supplied or null, the zero value
  // of the first element of the array is used. If the array is empty, @init
  // is returned.
  accumulate(arr, init=null)::
    local zero = if std.length(arr) > 0 then values.zero(arr[0]) else null;
    std.foldl(op.add, arr, array.coalesce([init, zero])),
}

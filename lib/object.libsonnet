{
  local object = self,
  local array = import 'array.libsonnet',

  // The nine functions {keys,vals,kv}{,All,Hidden} convert an object into
  // an array of values.
  //
  // The `keys` functions return an array of values that are the names of the fields (keys).
  // The `vals` functions return an array of values that are the values of the fields.
  // The `kv` functions return an array of values that are arrays of `[key, value]` pairs.
  //
  // The functions without a suffix return the visible (non-hidden) fields.
  // The functions with the `Hidden` suffix return the hidden fields.
  // The functions with the `All` suffix return all the fields.

  keys(obj):: std.objectFields(obj),
  keysAll(obj):: std.objectFieldsAll(obj),
  keysHidden(obj):: std.setDiff(object.fieldsAll(obj), object.fields(obj)),

  vals(obj):: __vals(obj, object.keys),
  valsAll(obj):: __vals(obj, object.keysAll),
  valsHidden(obj):: __vals(obj, object.keysHidden),
  local __vals(obj, selector) = [obj[k] for k in selector(obj)],

  kv(obj):: __kv(obj, object.keys),
  kvAll(obj):: __kv(obj, object.keysAll),
  kvHidden(obj):: __kv(obj, object.keysHidden),
  local __kv(obj, selector) = [[k, obj[k]] for k in selector(obj)],

  // make returns an object from an array @arr of `[key, value]` pairs. The
  // fields in the returned object are all visible (non-hidden).
  make(arr):: { [kv[0]]: kv[1] for kv in arr },

  // makeHidden returns an object from an array @arr of `[key, value]` pairs.
  // The fields in the returned object are all hidden.
  makeHidden(arr):: array.accumulate([{ [kv[0]]:: kv[1] } for kv in arr], {}),

  // transform calls @fn on each field of @obj passing it the key and the
  // value of each field and using the result of @fn to build the result of
  // transform. If @fn returns null, the field is removed from the result of
  // transform. Otherwise @fn should return a `[key, value]` pair which is
  // put into the result of transform as a field. The visibility of the field
  // in the result object is the same as the visibility of the field passed
  // to @fn.
  transform(fn, obj)::
    local visible = array.collapse(std.map(fn, object.kv(obj)));
    local hidden = array.collapse(std.map(fn, object.kvHidden(obj)));
    object.make(visible) + object.makeHidden(hidden),

  // filter calls @fn on each field of @obj passing it the key and the value
  // of each field and returns an object with the fields for which @fn
  // returned true.
  filter(fn, obj)::
    local map_fn(k, v) = if fn(k, v) then [k, v] else null;
    object.transform(map_fn, obj),

  // removeField returns @obj with @field removed from it.
  removeField(obj, field):: object.removeFields(obj, [field]),

  // removeFields returns @obj with each field listed in @fields removed
  // from it.
  removeFields(obj, fields)::
    local remove = std.set(fields);
    local keep_fn(k, v) = !std.setMember(k, remove);
    object.filter(keep_fn, obj),

  // invert swaps keys and values in @obj. The values in the output object is
  // a list of keys of visible fields from @obj. If any of the values in @obj
  // cannot be keys in an object (i.e. non-string values), jsonnet will
  // produce an evaluation error.
  // `{ a: 'foo', b: 'bar', c: 'foo' }` -> `{ bar: ['b'], foo: ['a', 'c'] }`
  invert(obj):: __invert(obj, object.keys),
  invertAll(obj):: _invert(obj, object.keysAll),
  local __invert(obj, selector) = array.accumulate([{ [obj[k]]+: [k] } for k in selector(obj)]),
}

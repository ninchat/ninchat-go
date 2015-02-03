"use strict";
(function($topLevelThis) {

Error.stackTraceLimit = Infinity;

var $global, $module;
if (typeof window !== "undefined") { /* web page */
  $global = window;
} else if (typeof self !== "undefined") { /* web worker */
  $global = self;
} else if (typeof global !== "undefined") { /* Node.js */
  $global = global;
  $global.require = require;
} else { /* others (e.g. Nashorn) */
  $global = $topLevelThis;
}

if ($global === undefined || $global.Array === undefined) {
  throw new Error("no global object found");
}
if (typeof module !== "undefined") {
  $module = module;
}

var $packages = {}, $reflect, $idCounter = 0;
var $keys = function(m) { return m ? Object.keys(m) : []; };
var $min = Math.min;
var $mod = function(x, y) { return x % y; };
var $parseInt = parseInt;
var $parseFloat = function(f) {
  if (f !== undefined && f !== null && f.constructor === Number) {
    return f;
  }
  return parseFloat(f);
};
var $flushConsole = function() {};

var $mapArray = function(array, f) {
  var newArray = new array.constructor(array.length), i;
  for (i = 0; i < array.length; i++) {
    newArray[i] = f(array[i]);
  }
  return newArray;
};

var $methodVal = function(recv, name) {
  var vals = recv.$methodVals || {};
  recv.$methodVals = vals; /* noop for primitives */
  var f = vals[name];
  if (f !== undefined) {
    return f;
  }
  var method = recv[name];
  f = function() {
    $stackDepthOffset--;
    try {
      return method.apply(recv, arguments);
    } finally {
      $stackDepthOffset++;
    }
  };
  vals[name] = f;
  return f;
};

var $methodExpr = function(method) {
  if (method.$expr === undefined) {
    method.$expr = function() {
      $stackDepthOffset--;
      try {
        return Function.call.apply(method, arguments);
      } finally {
        $stackDepthOffset++;
      }
    };
  }
  return method.$expr;
};

var $subslice = function(slice, low, high, max) {
  if (low < 0 || high < low || max < high || high > slice.$capacity || max > slice.$capacity) {
    $throwRuntimeError("slice bounds out of range");
  }
  var s = new slice.constructor(slice.$array);
  s.$offset = slice.$offset + low;
  s.$length = slice.$length - low;
  s.$capacity = slice.$capacity - low;
  if (high !== undefined) {
    s.$length = high - low;
  }
  if (max !== undefined) {
    s.$capacity = max - low;
  }
  return s;
};

var $sliceToArray = function(slice) {
  if (slice.$length === 0) {
    return [];
  }
  if (slice.$array.constructor !== Array) {
    return slice.$array.subarray(slice.$offset, slice.$offset + slice.$length);
  }
  return slice.$array.slice(slice.$offset, slice.$offset + slice.$length);
};

var $decodeRune = function(str, pos) {
  var c0 = str.charCodeAt(pos);

  if (c0 < 0x80) {
    return [c0, 1];
  }

  if (c0 !== c0 || c0 < 0xC0) {
    return [0xFFFD, 1];
  }

  var c1 = str.charCodeAt(pos + 1);
  if (c1 !== c1 || c1 < 0x80 || 0xC0 <= c1) {
    return [0xFFFD, 1];
  }

  if (c0 < 0xE0) {
    var r = (c0 & 0x1F) << 6 | (c1 & 0x3F);
    if (r <= 0x7F) {
      return [0xFFFD, 1];
    }
    return [r, 2];
  }

  var c2 = str.charCodeAt(pos + 2);
  if (c2 !== c2 || c2 < 0x80 || 0xC0 <= c2) {
    return [0xFFFD, 1];
  }

  if (c0 < 0xF0) {
    var r = (c0 & 0x0F) << 12 | (c1 & 0x3F) << 6 | (c2 & 0x3F);
    if (r <= 0x7FF) {
      return [0xFFFD, 1];
    }
    if (0xD800 <= r && r <= 0xDFFF) {
      return [0xFFFD, 1];
    }
    return [r, 3];
  }

  var c3 = str.charCodeAt(pos + 3);
  if (c3 !== c3 || c3 < 0x80 || 0xC0 <= c3) {
    return [0xFFFD, 1];
  }

  if (c0 < 0xF8) {
    var r = (c0 & 0x07) << 18 | (c1 & 0x3F) << 12 | (c2 & 0x3F) << 6 | (c3 & 0x3F);
    if (r <= 0xFFFF || 0x10FFFF < r) {
      return [0xFFFD, 1];
    }
    return [r, 4];
  }

  return [0xFFFD, 1];
};

var $encodeRune = function(r) {
  if (r < 0 || r > 0x10FFFF || (0xD800 <= r && r <= 0xDFFF)) {
    r = 0xFFFD;
  }
  if (r <= 0x7F) {
    return String.fromCharCode(r);
  }
  if (r <= 0x7FF) {
    return String.fromCharCode(0xC0 | r >> 6, 0x80 | (r & 0x3F));
  }
  if (r <= 0xFFFF) {
    return String.fromCharCode(0xE0 | r >> 12, 0x80 | (r >> 6 & 0x3F), 0x80 | (r & 0x3F));
  }
  return String.fromCharCode(0xF0 | r >> 18, 0x80 | (r >> 12 & 0x3F), 0x80 | (r >> 6 & 0x3F), 0x80 | (r & 0x3F));
};

var $stringToBytes = function(str) {
  var array = new Uint8Array(str.length), i;
  for (i = 0; i < str.length; i++) {
    array[i] = str.charCodeAt(i);
  }
  return array;
};

var $bytesToString = function(slice) {
  if (slice.$length === 0) {
    return "";
  }
  var str = "", i;
  for (i = 0; i < slice.$length; i += 10000) {
    str += String.fromCharCode.apply(null, slice.$array.subarray(slice.$offset + i, slice.$offset + Math.min(slice.$length, i + 10000)));
  }
  return str;
};

var $stringToRunes = function(str) {
  var array = new Int32Array(str.length);
  var rune, i, j = 0;
  for (i = 0; i < str.length; i += rune[1], j++) {
    rune = $decodeRune(str, i);
    array[j] = rune[0];
  }
  return array.subarray(0, j);
};

var $runesToString = function(slice) {
  if (slice.$length === 0) {
    return "";
  }
  var str = "", i;
  for (i = 0; i < slice.$length; i++) {
    str += $encodeRune(slice.$array[slice.$offset + i]);
  }
  return str;
};

var $copyString = function(dst, src) {
  var n = Math.min(src.length, dst.$length), i;
  for (i = 0; i < n; i++) {
    dst.$array[dst.$offset + i] = src.charCodeAt(i);
  }
  return n;
};

var $copySlice = function(dst, src) {
  var n = Math.min(src.$length, dst.$length), i;
  $internalCopy(dst.$array, src.$array, dst.$offset, src.$offset, n, dst.constructor.elem);
  return n;
};

var $copy = function(dst, src, type) {
  switch (type.kind) {
  case $kindArray:
    $internalCopy(dst, src, 0, 0, src.length, type.elem);
    break;
  case $kindStruct:
    for (var i = 0; i < type.fields.length; i++) {
      var field = type.fields[i];
      var fieldName = field[0];
      var fieldType = field[3];
      switch (fieldType.kind) {
      case $kindArray:
      case $kindStruct:
        $copy(dst[fieldName], src[fieldName], fieldType);
        continue;
      default:
        dst[fieldName] = src[fieldName];
        continue;
      }
    }
    break;
  }
};

var $internalCopy = function(dst, src, dstOffset, srcOffset, n, elem) {
  if (n === 0 || (dst === src && dstOffset === srcOffset)) {
    return;
  }

  if (src.subarray) {
    dst.set(src.subarray(srcOffset, srcOffset + n), dstOffset);
    return;
  }

  switch (elem.kind) {
  case $kindArray:
  case $kindStruct:
    if (dst === src && dstOffset > srcOffset) {
      for (var i = n - 1; i >= 0; i--) {
        $copy(dst[dstOffset + i], src[srcOffset + i], elem);
      }
      return;
    }
    for (var i = 0; i < n; i++) {
      $copy(dst[dstOffset + i], src[srcOffset + i], elem);
    }
    return;
  }

  if (dst === src && dstOffset > srcOffset) {
    for (var i = n - 1; i >= 0; i--) {
      dst[dstOffset + i] = src[srcOffset + i];
    }
    return;
  }
  for (var i = 0; i < n; i++) {
    dst[dstOffset + i] = src[srcOffset + i];
  }
};

var $clone = function(src, type) {
  var clone = type.zero();
  $copy(clone, src, type);
  return clone;
};

var $pointerOfStructConversion = function(obj, type) {
  if(obj.$proxies === undefined) {
    obj.$proxies = {};
    obj.$proxies[obj.constructor.string] = obj;
  }
  var proxy = obj.$proxies[type.string];
  if (proxy === undefined) {
    var properties = {};
    for (var i = 0; i < type.Struct.fields.length; i++) {
      (function(fieldName) {
        properties[fieldName] = {
          get: function() { return obj[fieldName]; },
          set: function(value) { obj[fieldName] = value; },
        };
      })(type.Struct.fields[i][0]);
    }
    proxy = Object.create(type.prototype, properties);
    proxy.$val = proxy;
    obj.$proxies[type.string] = proxy;
    proxy.$proxies = obj.$proxies;
  }
  return proxy;
};

var $append = function(slice) {
  return $internalAppend(slice, arguments, 1, arguments.length - 1);
};

var $appendSlice = function(slice, toAppend) {
  return $internalAppend(slice, toAppend.$array, toAppend.$offset, toAppend.$length);
};

var $internalAppend = function(slice, array, offset, length) {
  if (length === 0) {
    return slice;
  }

  var newArray = slice.$array;
  var newOffset = slice.$offset;
  var newLength = slice.$length + length;
  var newCapacity = slice.$capacity;

  if (newLength > newCapacity) {
    newOffset = 0;
    newCapacity = Math.max(newLength, slice.$capacity < 1024 ? slice.$capacity * 2 : Math.floor(slice.$capacity * 5 / 4));

    if (slice.$array.constructor === Array) {
      newArray = slice.$array.slice(slice.$offset, slice.$offset + slice.$length);
      newArray.length = newCapacity;
      var zero = slice.constructor.elem.zero, i;
      for (i = slice.$length; i < newCapacity; i++) {
        newArray[i] = zero();
      }
    } else {
      newArray = new slice.$array.constructor(newCapacity);
      newArray.set(slice.$array.subarray(slice.$offset, slice.$offset + slice.$length));
    }
  }

  $internalCopy(newArray, array, newOffset + slice.$length, offset, length, slice.constructor.elem);

  var newSlice = new slice.constructor(newArray);
  newSlice.$offset = newOffset;
  newSlice.$length = newLength;
  newSlice.$capacity = newCapacity;
  return newSlice;
};

var $equal = function(a, b, type) {
  if (a === b) {
    return true;
  }
  var i;
  switch (type.kind) {
  case $kindFloat32:
    return $float32IsEqual(a, b);
  case $kindComplex64:
    return $float32IsEqual(a.$real, b.$real) && $float32IsEqual(a.$imag, b.$imag);
  case $kindComplex128:
    return a.$real === b.$real && a.$imag === b.$imag;
  case $kindInt64:
  case $kindUint64:
    return a.$high === b.$high && a.$low === b.$low;
  case $kindPtr:
    if (a.constructor.Struct) {
      return false;
    }
    return $pointerIsEqual(a, b);
  case $kindArray:
    if (a.length != b.length) {
      return false;
    }
    var i;
    for (i = 0; i < a.length; i++) {
      if (!$equal(a[i], b[i], type.elem)) {
        return false;
      }
    }
    return true;
  case $kindStruct:
    for (i = 0; i < type.fields.length; i++) {
      var field = type.fields[i];
      var name = field[0];
      if (!$equal(a[name], b[name], field[3])) {
        return false;
      }
    }
    return true;
  default:
    return false;
  }
};

var $interfaceIsEqual = function(a, b) {
  if (a === null || b === null || a === undefined || b === undefined || a.constructor !== b.constructor) {
    return a === b;
  }
  switch (a.constructor.kind) {
  case $kindFunc:
  case $kindMap:
  case $kindSlice:
  case $kindStruct:
    $throwRuntimeError("comparing uncomparable type " + a.constructor.string);
  case undefined: /* js.Object */
    return a === b;
  default:
    return $equal(a.$val, b.$val, a.constructor);
  }
};

var $float32IsEqual = function(a, b) {
  if (a === b) {
    return true;
  }
  if (a === 1/0 || b === 1/0 || a === -1/0 || b === -1/0 || a !== a || b !== b) {
    return false;
  }
  var math = $packages["math"];
  return math !== undefined && math.Float32bits(a) === math.Float32bits(b);
};

var $pointerIsEqual = function(a, b) {
  if (a === b) {
    return true;
  }
  if (a.$get === $throwNilPointerError || b.$get === $throwNilPointerError) {
    return a.$get === $throwNilPointerError && b.$get === $throwNilPointerError;
  }
  var va = a.$get();
  var vb = b.$get();
  if (va !== vb) {
    return false;
  }
  var dummy = va + 1;
  a.$set(dummy);
  var equal = b.$get() === dummy;
  a.$set(va);
  return equal;
};

var $kindBool = 1;
var $kindInt = 2;
var $kindInt8 = 3;
var $kindInt16 = 4;
var $kindInt32 = 5;
var $kindInt64 = 6;
var $kindUint = 7;
var $kindUint8 = 8;
var $kindUint16 = 9;
var $kindUint32 = 10;
var $kindUint64 = 11;
var $kindUintptr = 12;
var $kindFloat32 = 13;
var $kindFloat64 = 14;
var $kindComplex64 = 15;
var $kindComplex128 = 16;
var $kindArray = 17;
var $kindChan = 18;
var $kindFunc = 19;
var $kindInterface = 20;
var $kindMap = 21;
var $kindPtr = 22;
var $kindSlice = 23;
var $kindString = 24;
var $kindStruct = 25;
var $kindUnsafePointer = 26;

var $newType = function(size, kind, string, name, pkgPath, constructor) {
  var typ;
  switch(kind) {
  case $kindBool:
  case $kindInt:
  case $kindInt8:
  case $kindInt16:
  case $kindInt32:
  case $kindUint:
  case $kindUint8:
  case $kindUint16:
  case $kindUint32:
  case $kindUintptr:
  case $kindString:
  case $kindUnsafePointer:
    typ = function(v) { this.$val = v; };
    typ.prototype.$key = function() { return string + "$" + this.$val; };
    break;

  case $kindFloat32:
  case $kindFloat64:
    typ = function(v) { this.$val = v; };
    typ.prototype.$key = function() { return string + "$" + $floatKey(this.$val); };
    break;

  case $kindInt64:
    typ = function(high, low) {
      this.$high = (high + Math.floor(Math.ceil(low) / 4294967296)) >> 0;
      this.$low = low >>> 0;
      this.$val = this;
    };
    typ.prototype.$key = function() { return string + "$" + this.$high + "$" + this.$low; };
    break;

  case $kindUint64:
    typ = function(high, low) {
      this.$high = (high + Math.floor(Math.ceil(low) / 4294967296)) >>> 0;
      this.$low = low >>> 0;
      this.$val = this;
    };
    typ.prototype.$key = function() { return string + "$" + this.$high + "$" + this.$low; };
    break;

  case $kindComplex64:
  case $kindComplex128:
    typ = function(real, imag) {
      this.$real = real;
      this.$imag = imag;
      this.$val = this;
    };
    typ.prototype.$key = function() { return string + "$" + this.$real + "$" + this.$imag; };
    break;

  case $kindArray:
    typ = function(v) { this.$val = v; };
    typ.Ptr = $newType(4, $kindPtr, "*" + string, "", "", function(array) {
      this.$get = function() { return array; };
      this.$set = function(v) { $copy(this, v, typ); };
      this.$val = array;
    });
    typ.init = function(elem, len) {
      typ.elem = elem;
      typ.len = len;
      typ.prototype.$key = function() {
        return string + "$" + Array.prototype.join.call($mapArray(this.$val, function(e) {
          var key = e.$key ? e.$key() : String(e);
          return key.replace(/\\/g, "\\\\").replace(/\$/g, "\\$");
        }), "$");
      };
      typ.extendReflectType = function(rt) {
        rt.arrayType = new $reflect.arrayType.Ptr(rt, elem.reflectType(), undefined, len);
      };
      typ.Ptr.init(typ);
      Object.defineProperty(typ.Ptr.nil, "nilCheck", { get: $throwNilPointerError });
    };
    break;

  case $kindChan:
    typ = function(capacity) {
      this.$val = this;
      this.$capacity = capacity;
      this.$buffer = [];
      this.$sendQueue = [];
      this.$recvQueue = [];
      this.$closed = false;
    };
    typ.prototype.$key = function() {
      if (this.$id === undefined) {
        $idCounter++;
        this.$id = $idCounter;
      }
      return String(this.$id);
    };
    typ.init = function(elem, sendOnly, recvOnly) {
      typ.elem = elem;
      typ.sendOnly = sendOnly;
      typ.recvOnly = recvOnly;
      typ.nil = new typ(0);
      typ.nil.$sendQueue = typ.nil.$recvQueue = { length: 0, push: function() {}, shift: function() { return undefined; }, indexOf: function() { return -1; } };
      typ.extendReflectType = function(rt) {
        rt.chanType = new $reflect.chanType.Ptr(rt, elem.reflectType(), sendOnly ? $reflect.SendDir : (recvOnly ? $reflect.RecvDir : $reflect.BothDir));
      };
    };
    break;

  case $kindFunc:
    typ = function(v) { this.$val = v; };
    typ.init = function(params, results, variadic) {
      typ.params = params;
      typ.results = results;
      typ.variadic = variadic;
      typ.extendReflectType = function(rt) {
        var typeSlice = ($sliceType($ptrType($reflect.rtype.Ptr)));
        rt.funcType = new $reflect.funcType.Ptr(rt, variadic, new typeSlice($mapArray(params, function(p) { return p.reflectType(); })), new typeSlice($mapArray(results, function(p) { return p.reflectType(); })));
      };
    };
    break;

  case $kindInterface:
    typ = { implementedBy: {}, missingMethodFor: {} };
    typ.init = function(methods) {
      typ.methods = methods;
      typ.extendReflectType = function(rt) {
        var imethods = $mapArray(methods, function(m) {
          return new $reflect.imethod.Ptr($newStringPtr(m[1]), $newStringPtr(m[2]), m[3].reflectType());
        });
        var methodSlice = ($sliceType($ptrType($reflect.imethod.Ptr)));
        rt.interfaceType = new $reflect.interfaceType.Ptr(rt, new methodSlice(imethods));
      };
    };
    break;

  case $kindMap:
    typ = function(v) { this.$val = v; };
    typ.init = function(key, elem) {
      typ.key = key;
      typ.elem = elem;
      typ.extendReflectType = function(rt) {
        rt.mapType = new $reflect.mapType.Ptr(rt, key.reflectType(), elem.reflectType(), undefined, undefined);
      };
    };
    break;

  case $kindPtr:
    typ = constructor || function(getter, setter, target) {
      this.$get = getter;
      this.$set = setter;
      this.$target = target;
      this.$val = this;
    };
    typ.prototype.$key = function() {
      if (this.$id === undefined) {
        $idCounter++;
        this.$id = $idCounter;
      }
      return String(this.$id);
    };
    typ.init = function(elem) {
      typ.nil = new typ($throwNilPointerError, $throwNilPointerError);
      typ.extendReflectType = function(rt) {
        rt.ptrType = new $reflect.ptrType.Ptr(rt, elem.reflectType());
      };
    };
    break;

  case $kindSlice:
    var nativeArray;
    typ = function(array) {
      if (array.constructor !== nativeArray) {
        array = new nativeArray(array);
      }
      this.$array = array;
      this.$offset = 0;
      this.$length = array.length;
      this.$capacity = array.length;
      this.$val = this;
    };
    typ.make = function(length, capacity) {
      capacity = capacity || length;
      var array = new nativeArray(capacity), i;
      if (nativeArray === Array) {
        for (i = 0; i < capacity; i++) {
          array[i] = typ.elem.zero();
        }
      }
      var slice = new typ(array);
      slice.$length = length;
      return slice;
    };
    typ.init = function(elem) {
      typ.elem = elem;
      nativeArray = $nativeArray(elem.kind);
      typ.nil = new typ([]);
      typ.extendReflectType = function(rt) {
        rt.sliceType = new $reflect.sliceType.Ptr(rt, elem.reflectType());
      };
    };
    break;

  case $kindStruct:
    typ = function(v) { this.$val = v; };
    typ.Ptr = $newType(4, $kindPtr, "*" + string, "", "", constructor);
    typ.Ptr.Struct = typ;
    typ.Ptr.prototype.$get = function() { return this; };
    typ.Ptr.prototype.$set = function(v) { $copy(this, v, typ); };
    typ.init = function(fields) {
      var i;
      typ.fields = fields;
      typ.prototype.$key = function() {
        var val = this.$val;
        return string + "$" + $mapArray(fields, function(field) {
          var e = val[field[0]];
          var key = e.$key ? e.$key() : String(e);
          return key.replace(/\\/g, "\\\\").replace(/\$/g, "\\$");
        }).join("$");
      };
      typ.Ptr.extendReflectType = function(rt) {
        rt.ptrType = new $reflect.ptrType.Ptr(rt, typ.reflectType());
      };
      /* nil value */
      var properties = {};
      for (i = 0; i < fields.length; i++) {
        properties[fields[i][0]] = { get: $throwNilPointerError, set: $throwNilPointerError };
      }
      typ.Ptr.nil = Object.create(constructor.prototype, properties);
      typ.Ptr.nil.$val = typ.Ptr.nil;
      /* methods for embedded fields */
      for (i = 0; i < typ.methods.length; i++) {
        var m = typ.methods[i];
        if (m[4] != -1) {
          (function(field, methodName) {
            typ.prototype[methodName] = function() {
              var v = this.$val[field[0]];
              return v[methodName].apply(v, arguments);
            };
          })(fields[m[4]], m[0]);
        }
      }
      for (i = 0; i < typ.Ptr.methods.length; i++) {
        var m = typ.Ptr.methods[i];
        if (m[4] != -1) {
          (function(field, methodName) {
            typ.Ptr.prototype[methodName] = function() {
              var v = this[field[0]];
              if (v.$val === undefined) {
                v = new field[3](v);
              }
              return v[methodName].apply(v, arguments);
            };
          })(fields[m[4]], m[0]);
        }
      }
      /* reflect type */
      typ.extendReflectType = function(rt) {
        var reflectFields = new Array(fields.length), i;
        for (i = 0; i < fields.length; i++) {
          var field = fields[i];
          reflectFields[i] = new $reflect.structField.Ptr($newStringPtr(field[1]), $newStringPtr(field[2]), field[3].reflectType(), $newStringPtr(field[4]), i);
        }
        rt.structType = new $reflect.structType.Ptr(rt, new ($sliceType($reflect.structField.Ptr))(reflectFields));
      };
    };
    break;

  default:
    $panic(new $String("invalid kind: " + kind));
  }

  switch(kind) {
  case $kindBool:
  case $kindMap:
    typ.zero = function() { return false; };
    break;

  case $kindInt:
  case $kindInt8:
  case $kindInt16:
  case $kindInt32:
  case $kindUint:
  case $kindUint8 :
  case $kindUint16:
  case $kindUint32:
  case $kindUintptr:
  case $kindUnsafePointer:
  case $kindFloat32:
  case $kindFloat64:
    typ.zero = function() { return 0; };
    break;

  case $kindString:
    typ.zero = function() { return ""; };
    break;

  case $kindInt64:
  case $kindUint64:
  case $kindComplex64:
  case $kindComplex128:
    var zero = new typ(0, 0);
    typ.zero = function() { return zero; };
    break;

  case $kindChan:
  case $kindPtr:
  case $kindSlice:
    typ.zero = function() { return typ.nil; };
    break;

  case $kindFunc:
    typ.zero = function() { return $throwNilPointerError; };
    break;

  case $kindInterface:
    typ.zero = function() { return $ifaceNil; };
    break;

  case $kindArray:
    typ.zero = function() {
      var arrayClass = $nativeArray(typ.elem.kind);
      if (arrayClass !== Array) {
        return new arrayClass(typ.len);
      }
      var array = new Array(typ.len), i;
      for (i = 0; i < typ.len; i++) {
        array[i] = typ.elem.zero();
      }
      return array;
    };
    break;

  case $kindStruct:
    typ.zero = function() { return new typ.Ptr(); };
    break;

  default:
    $panic(new $String("invalid kind: " + kind));
  }

  typ.kind = kind;
  typ.string = string;
  typ.typeName = name;
  typ.pkgPath = pkgPath;
  typ.methods = [];
  var rt = null;
  typ.reflectType = function() {
    if (rt === null) {
      rt = new $reflect.rtype.Ptr(size, 0, 0, 0, 0, kind, undefined, undefined, $newStringPtr(string), undefined, undefined);
      rt.jsType = typ;

      var methods = [];
      if (typ.methods !== undefined) {
        var i;
        for (i = 0; i < typ.methods.length; i++) {
          var m = typ.methods[i];
          var t = m[3];
          methods.push(new $reflect.method.Ptr($newStringPtr(m[1]), $newStringPtr(m[2]), t.reflectType(), $funcType([typ].concat(t.params), t.results, t.variadic).reflectType(), undefined, undefined));
        }
      }
      if (name !== "" || methods.length !== 0) {
        var methodSlice = ($sliceType($ptrType($reflect.method.Ptr)));
        rt.uncommonType = new $reflect.uncommonType.Ptr($newStringPtr(name), $newStringPtr(pkgPath), new methodSlice(methods));
        rt.uncommonType.jsType = typ;
      }

      if (typ.extendReflectType !== undefined) {
        typ.extendReflectType(rt);
      }
    }
    return rt;
  };
  return typ;
};

var $Bool          = $newType( 1, $kindBool,          "bool",           "bool",       "", null);
var $Int           = $newType( 4, $kindInt,           "int",            "int",        "", null);
var $Int8          = $newType( 1, $kindInt8,          "int8",           "int8",       "", null);
var $Int16         = $newType( 2, $kindInt16,         "int16",          "int16",      "", null);
var $Int32         = $newType( 4, $kindInt32,         "int32",          "int32",      "", null);
var $Int64         = $newType( 8, $kindInt64,         "int64",          "int64",      "", null);
var $Uint          = $newType( 4, $kindUint,          "uint",           "uint",       "", null);
var $Uint8         = $newType( 1, $kindUint8,         "uint8",          "uint8",      "", null);
var $Uint16        = $newType( 2, $kindUint16,        "uint16",         "uint16",     "", null);
var $Uint32        = $newType( 4, $kindUint32,        "uint32",         "uint32",     "", null);
var $Uint64        = $newType( 8, $kindUint64,        "uint64",         "uint64",     "", null);
var $Uintptr       = $newType( 4, $kindUintptr,       "uintptr",        "uintptr",    "", null);
var $Float32       = $newType( 4, $kindFloat32,       "float32",        "float32",    "", null);
var $Float64       = $newType( 8, $kindFloat64,       "float64",        "float64",    "", null);
var $Complex64     = $newType( 8, $kindComplex64,     "complex64",      "complex64",  "", null);
var $Complex128    = $newType(16, $kindComplex128,    "complex128",     "complex128", "", null);
var $String        = $newType( 8, $kindString,        "string",         "string",     "", null);
var $UnsafePointer = $newType( 4, $kindUnsafePointer, "unsafe.Pointer", "Pointer",    "", null);

var $nativeArray = function(elemKind) {
  switch (elemKind) {
  case $kindInt:
    return Int32Array;
  case $kindInt8:
    return Int8Array;
  case $kindInt16:
    return Int16Array;
  case $kindInt32:
    return Int32Array;
  case $kindUint:
    return Uint32Array;
  case $kindUint8:
    return Uint8Array;
  case $kindUint16:
    return Uint16Array;
  case $kindUint32:
    return Uint32Array;
  case $kindUintptr:
    return Uint32Array;
  case $kindFloat32:
    return Float32Array;
  case $kindFloat64:
    return Float64Array;
  default:
    return Array;
  }
};
var $toNativeArray = function(elemKind, array) {
  var nativeArray = $nativeArray(elemKind);
  if (nativeArray === Array) {
    return array;
  }
  return new nativeArray(array);
};
var $arrayTypes = {};
var $arrayType = function(elem, len) {
  var string = "[" + len + "]" + elem.string;
  var typ = $arrayTypes[string];
  if (typ === undefined) {
    typ = $newType(12, $kindArray, string, "", "", null);
    typ.init(elem, len);
    $arrayTypes[string] = typ;
  }
  return typ;
};

var $chanType = function(elem, sendOnly, recvOnly) {
  var string = (recvOnly ? "<-" : "") + "chan" + (sendOnly ? "<- " : " ") + elem.string;
  var field = sendOnly ? "SendChan" : (recvOnly ? "RecvChan" : "Chan");
  var typ = elem[field];
  if (typ === undefined) {
    typ = $newType(4, $kindChan, string, "", "", null);
    typ.init(elem, sendOnly, recvOnly);
    elem[field] = typ;
  }
  return typ;
};

var $funcTypes = {};
var $funcType = function(params, results, variadic) {
  var paramTypes = $mapArray(params, function(p) { return p.string; });
  if (variadic) {
    paramTypes[paramTypes.length - 1] = "..." + paramTypes[paramTypes.length - 1].substr(2);
  }
  var string = "func(" + paramTypes.join(", ") + ")";
  if (results.length === 1) {
    string += " " + results[0].string;
  } else if (results.length > 1) {
    string += " (" + $mapArray(results, function(r) { return r.string; }).join(", ") + ")";
  }
  var typ = $funcTypes[string];
  if (typ === undefined) {
    typ = $newType(4, $kindFunc, string, "", "", null);
    typ.init(params, results, variadic);
    $funcTypes[string] = typ;
  }
  return typ;
};

var $interfaceTypes = {};
var $interfaceType = function(methods) {
  var string = "interface {}";
  if (methods.length !== 0) {
    string = "interface { " + $mapArray(methods, function(m) {
      return (m[2] !== "" ? m[2] + "." : "") + m[1] + m[3].string.substr(4);
    }).join("; ") + " }";
  }
  var typ = $interfaceTypes[string];
  if (typ === undefined) {
    typ = $newType(8, $kindInterface, string, "", "", null);
    typ.init(methods);
    $interfaceTypes[string] = typ;
  }
  return typ;
};
var $emptyInterface = $interfaceType([]);
var $ifaceNil = { $key: function() { return "nil"; } };
var $error = $newType(8, $kindInterface, "error", "error", "", null);
$error.init([["Error", "Error", "", $funcType([], [$String], false)]]);

var $Map = function() {};
(function() {
  var names = Object.getOwnPropertyNames(Object.prototype), i;
  for (i = 0; i < names.length; i++) {
    $Map.prototype[names[i]] = undefined;
  }
})();
var $mapTypes = {};
var $mapType = function(key, elem) {
  var string = "map[" + key.string + "]" + elem.string;
  var typ = $mapTypes[string];
  if (typ === undefined) {
    typ = $newType(4, $kindMap, string, "", "", null);
    typ.init(key, elem);
    $mapTypes[string] = typ;
  }
  return typ;
};


var $throwNilPointerError = function() { $throwRuntimeError("invalid memory address or nil pointer dereference"); };
var $ptrType = function(elem) {
  var typ = elem.Ptr;
  if (typ === undefined) {
    typ = $newType(4, $kindPtr, "*" + elem.string, "", "", null);
    typ.init(elem);
    elem.Ptr = typ;
  }
  return typ;
};

var $stringPtrMap = new $Map();
var $newStringPtr = function(str) {
  if (str === undefined || str === "") {
    return $ptrType($String).nil;
  }
  var ptr = $stringPtrMap[str];
  if (ptr === undefined) {
    ptr = new ($ptrType($String))(function() { return str; }, function(v) { str = v; });
    $stringPtrMap[str] = ptr;
  }
  return ptr;
};

var $newDataPointer = function(data, constructor) {
  if (constructor.Struct) {
    return data;
  }
  return new constructor(function() { return data; }, function(v) { data = v; });
};

var $sliceType = function(elem) {
  var typ = elem.Slice;
  if (typ === undefined) {
    typ = $newType(12, $kindSlice, "[]" + elem.string, "", "", null);
    typ.init(elem);
    elem.Slice = typ;
  }
  return typ;
};

var $structTypes = {};
var $structType = function(fields) {
  var string = "struct { " + $mapArray(fields, function(f) {
    return f[1] + " " + f[3].string + (f[4] !== "" ? (" \"" + f[4].replace(/\\/g, "\\\\").replace(/"/g, "\\\"") + "\"") : "");
  }).join("; ") + " }";
  if (fields.length === 0) {
    string = "struct {}";
  }
  var typ = $structTypes[string];
  if (typ === undefined) {
    typ = $newType(0, $kindStruct, string, "", "", function() {
      this.$val = this;
      var i;
      for (i = 0; i < fields.length; i++) {
        var field = fields[i];
        var arg = arguments[i];
        this[field[0]] = arg !== undefined ? arg : field[3].zero();
      }
    });
    /* collect methods for anonymous fields */
    var i, j;
    for (i = 0; i < fields.length; i++) {
      var field = fields[i];
      if (field[1] === "") {
        var methods = field[3].methods;
        for (j = 0; j < methods.length; j++) {
          var m = methods[j].slice(0, 6).concat([i]);
          typ.methods.push(m);
          typ.Ptr.methods.push(m);
        }
        if (field[3].kind === $kindStruct) {
          var methods = field[3].Ptr.methods;
          for (j = 0; j < methods.length; j++) {
            typ.Ptr.methods.push(methods[j].slice(0, 6).concat([i]));
          }
        }
      }
    }
    typ.init(fields);
    $structTypes[string] = typ;
  }
  return typ;
};

var $assertType = function(value, type, returnTuple) {
  var isInterface = (type.kind === $kindInterface), ok, missingMethod = "";
  if (value === $ifaceNil) {
    ok = false;
  } else if (!isInterface) {
    ok = value.constructor === type;
  } else if (type.string === "js.Object") {
    ok = true;
  } else {
    var valueTypeString = value.constructor.string;
    ok = type.implementedBy[valueTypeString];
    if (ok === undefined) {
      ok = true;
      var valueMethods = value.constructor.methods;
      var typeMethods = type.methods;
      for (var i = 0; i < typeMethods.length; i++) {
        var tm = typeMethods[i];
        var found = false;
        for (var j = 0; j < valueMethods.length; j++) {
          var vm = valueMethods[j];
          if (vm[1] === tm[1] && vm[2] === tm[2] && vm[3] === tm[3]) {
            found = true;
            break;
          }
        }
        if (!found) {
          ok = false;
          type.missingMethodFor[valueTypeString] = tm[1];
          break;
        }
      }
      type.implementedBy[valueTypeString] = ok;
    }
    if (!ok) {
      missingMethod = type.missingMethodFor[valueTypeString];
    }
  }

  if (!ok) {
    if (returnTuple) {
      return [type.zero(), false];
    }
    $panic(new $packages["runtime"].TypeAssertionError.Ptr("", (value === $ifaceNil ? "" : value.constructor.string), type.string, missingMethod));
  }

  if (!isInterface) {
    value = value.$val;
  }
  return returnTuple ? [value, true] : value;
};

var $coerceFloat32 = function(f) {
  var math = $packages["math"];
  if (math === undefined) {
    return f;
  }
  return math.Float32frombits(math.Float32bits(f));
};

var $floatKey = function(f) {
  if (f !== f) {
    $idCounter++;
    return "NaN$" + $idCounter;
  }
  return String(f);
};

var $flatten64 = function(x) {
  return x.$high * 4294967296 + x.$low;
};

var $shiftLeft64 = function(x, y) {
  if (y === 0) {
    return x;
  }
  if (y < 32) {
    return new x.constructor(x.$high << y | x.$low >>> (32 - y), (x.$low << y) >>> 0);
  }
  if (y < 64) {
    return new x.constructor(x.$low << (y - 32), 0);
  }
  return new x.constructor(0, 0);
};

var $shiftRightInt64 = function(x, y) {
  if (y === 0) {
    return x;
  }
  if (y < 32) {
    return new x.constructor(x.$high >> y, (x.$low >>> y | x.$high << (32 - y)) >>> 0);
  }
  if (y < 64) {
    return new x.constructor(x.$high >> 31, (x.$high >> (y - 32)) >>> 0);
  }
  if (x.$high < 0) {
    return new x.constructor(-1, 4294967295);
  }
  return new x.constructor(0, 0);
};

var $shiftRightUint64 = function(x, y) {
  if (y === 0) {
    return x;
  }
  if (y < 32) {
    return new x.constructor(x.$high >>> y, (x.$low >>> y | x.$high << (32 - y)) >>> 0);
  }
  if (y < 64) {
    return new x.constructor(0, x.$high >>> (y - 32));
  }
  return new x.constructor(0, 0);
};

var $mul64 = function(x, y) {
  var high = 0, low = 0, i;
  if ((y.$low & 1) !== 0) {
    high = x.$high;
    low = x.$low;
  }
  for (i = 1; i < 32; i++) {
    if ((y.$low & 1<<i) !== 0) {
      high += x.$high << i | x.$low >>> (32 - i);
      low += (x.$low << i) >>> 0;
    }
  }
  for (i = 0; i < 32; i++) {
    if ((y.$high & 1<<i) !== 0) {
      high += x.$low << i;
    }
  }
  return new x.constructor(high, low);
};

var $div64 = function(x, y, returnRemainder) {
  if (y.$high === 0 && y.$low === 0) {
    $throwRuntimeError("integer divide by zero");
  }

  var s = 1;
  var rs = 1;

  var xHigh = x.$high;
  var xLow = x.$low;
  if (xHigh < 0) {
    s = -1;
    rs = -1;
    xHigh = -xHigh;
    if (xLow !== 0) {
      xHigh--;
      xLow = 4294967296 - xLow;
    }
  }

  var yHigh = y.$high;
  var yLow = y.$low;
  if (y.$high < 0) {
    s *= -1;
    yHigh = -yHigh;
    if (yLow !== 0) {
      yHigh--;
      yLow = 4294967296 - yLow;
    }
  }

  var high = 0, low = 0, n = 0, i;
  while (yHigh < 2147483648 && ((xHigh > yHigh) || (xHigh === yHigh && xLow > yLow))) {
    yHigh = (yHigh << 1 | yLow >>> 31) >>> 0;
    yLow = (yLow << 1) >>> 0;
    n++;
  }
  for (i = 0; i <= n; i++) {
    high = high << 1 | low >>> 31;
    low = (low << 1) >>> 0;
    if ((xHigh > yHigh) || (xHigh === yHigh && xLow >= yLow)) {
      xHigh = xHigh - yHigh;
      xLow = xLow - yLow;
      if (xLow < 0) {
        xHigh--;
        xLow += 4294967296;
      }
      low++;
      if (low === 4294967296) {
        high++;
        low = 0;
      }
    }
    yLow = (yLow >>> 1 | yHigh << (32 - 1)) >>> 0;
    yHigh = yHigh >>> 1;
  }

  if (returnRemainder) {
    return new x.constructor(xHigh * rs, xLow * rs);
  }
  return new x.constructor(high * s, low * s);
};

var $divComplex = function(n, d) {
  var ninf = n.$real === 1/0 || n.$real === -1/0 || n.$imag === 1/0 || n.$imag === -1/0;
  var dinf = d.$real === 1/0 || d.$real === -1/0 || d.$imag === 1/0 || d.$imag === -1/0;
  var nnan = !ninf && (n.$real !== n.$real || n.$imag !== n.$imag);
  var dnan = !dinf && (d.$real !== d.$real || d.$imag !== d.$imag);
  if(nnan || dnan) {
    return new n.constructor(0/0, 0/0);
  }
  if (ninf && !dinf) {
    return new n.constructor(1/0, 1/0);
  }
  if (!ninf && dinf) {
    return new n.constructor(0, 0);
  }
  if (d.$real === 0 && d.$imag === 0) {
    if (n.$real === 0 && n.$imag === 0) {
      return new n.constructor(0/0, 0/0);
    }
    return new n.constructor(1/0, 1/0);
  }
  var a = Math.abs(d.$real);
  var b = Math.abs(d.$imag);
  if (a <= b) {
    var ratio = d.$real / d.$imag;
    var denom = d.$real * ratio + d.$imag;
    return new n.constructor((n.$real * ratio + n.$imag) / denom, (n.$imag * ratio - n.$real) / denom);
  }
  var ratio = d.$imag / d.$real;
  var denom = d.$imag * ratio + d.$real;
  return new n.constructor((n.$imag * ratio + n.$real) / denom, (n.$imag - n.$real * ratio) / denom);
};

var $stackDepthOffset = 0;
var $getStackDepth = function() {
  var err = new Error();
  if (err.stack === undefined) {
    return undefined;
  }
  return $stackDepthOffset + err.stack.split("\n").length;
};

var $deferFrames = [], $skippedDeferFrames = 0, $jumpToDefer = false, $panicStackDepth = null, $panicValue;
var $callDeferred = function(deferred, jsErr) {
  if ($skippedDeferFrames !== 0) {
    $skippedDeferFrames--;
    throw jsErr;
  }
  if ($jumpToDefer) {
    $jumpToDefer = false;
    throw jsErr;
  }
  if (jsErr) {
    var newErr = null;
    try {
      $deferFrames.push(deferred);
      $panic(new $packages["github.com/gopherjs/gopherjs/js"].Error.Ptr(jsErr));
    } catch (err) {
      newErr = err;
    }
    $deferFrames.pop();
    $callDeferred(deferred, newErr);
    return;
  }

  $stackDepthOffset--;
  var outerPanicStackDepth = $panicStackDepth;
  var outerPanicValue = $panicValue;

  var localPanicValue = $curGoroutine.panicStack.pop();
  if (localPanicValue !== undefined) {
    $panicStackDepth = $getStackDepth();
    $panicValue = localPanicValue;
  }

  var call, localSkippedDeferFrames = 0;
  try {
    while (true) {
      if (deferred === null) {
        deferred = $deferFrames[$deferFrames.length - 1 - localSkippedDeferFrames];
        if (deferred === undefined) {
          var msg;
          if (localPanicValue.constructor === $String) {
            msg = localPanicValue.$val;
          } else if (localPanicValue.Error !== undefined) {
            msg = localPanicValue.Error();
          } else if (localPanicValue.String !== undefined) {
            msg = localPanicValue.String();
          } else {
            msg = localPanicValue;
          }
          var e = new Error(msg);
          if (localPanicValue.Stack !== undefined) {
            e.stack = localPanicValue.Stack();
            e.stack = msg + e.stack.substr(e.stack.indexOf("\n"));
          }
          throw e;
        }
      }
      var call = deferred.pop();
      if (call === undefined) {
        if (localPanicValue !== undefined) {
          localSkippedDeferFrames++;
          deferred = null;
          continue;
        }
        return;
      }
      var r = call[0].apply(undefined, call[1]);
      if (r && r.$blocking) {
        deferred.push([r, []]);
      }

      if (localPanicValue !== undefined && $panicStackDepth === null) {
        throw null; /* error was recovered */
      }
    }
  } finally {
    $skippedDeferFrames += localSkippedDeferFrames;
    if ($curGoroutine.asleep) {
      deferred.push(call);
      $jumpToDefer = true;
    }
    if (localPanicValue !== undefined) {
      if ($panicStackDepth !== null) {
        $curGoroutine.panicStack.push(localPanicValue);
      }
      $panicStackDepth = outerPanicStackDepth;
      $panicValue = outerPanicValue;
    }
    $stackDepthOffset++;
  }
};

var $panic = function(value) {
  $curGoroutine.panicStack.push(value);
  $callDeferred(null, null);
};
var $recover = function() {
  if ($panicStackDepth === null || ($panicStackDepth !== undefined && $panicStackDepth !== $getStackDepth() - 2)) {
    return $ifaceNil;
  }
  $panicStackDepth = null;
  return $panicValue;
};
var $throw = function(err) { throw err; };
var $throwRuntimeError; /* set by package "runtime" */

var $BLOCKING = new Object();
var $nonblockingCall = function() {
  $panic(new $packages["runtime"].NotSupportedError.Ptr("non-blocking call to blocking function, see https://github.com/gopherjs/gopherjs#goroutines"));
};

var $dummyGoroutine = { asleep: false, exit: false, panicStack: [] };
var $curGoroutine = $dummyGoroutine, $totalGoroutines = 0, $awakeGoroutines = 0, $checkForDeadlock = true;
var $go = function(fun, args, direct) {
  $totalGoroutines++;
  $awakeGoroutines++;
  args.push($BLOCKING);
  var goroutine = function() {
    var rescheduled = false;
    try {
      $curGoroutine = goroutine;
      $skippedDeferFrames = 0;
      $jumpToDefer = false;
      var r = fun.apply(undefined, args);
      if (r && r.$blocking) {
        fun = r;
        args = [];
        $schedule(goroutine, direct);
        rescheduled = true;
        return;
      }
      goroutine.exit = true;
    } catch (err) {
      if (!$curGoroutine.asleep) {
        goroutine.exit = true;
        throw err;
      }
    } finally {
      $curGoroutine = $dummyGoroutine;
      if (goroutine.exit && !rescheduled) { /* also set by runtime.Goexit() */
        $totalGoroutines--;
        goroutine.asleep = true;
      }
      if (goroutine.asleep && !rescheduled) {
        $awakeGoroutines--;
        if ($awakeGoroutines === 0 && $totalGoroutines !== 0 && $checkForDeadlock) {
          console.error("fatal error: all goroutines are asleep - deadlock!");
        }
      }
    }
  };
  goroutine.asleep = false;
  goroutine.exit = false;
  goroutine.panicStack = [];
  $schedule(goroutine, direct);
};

var $scheduled = [], $schedulerLoopActive = false;
var $schedule = function(goroutine, direct) {
  if (goroutine.asleep) {
    goroutine.asleep = false;
    $awakeGoroutines++;
  }

  if (direct) {
    goroutine();
    return;
  }

  $scheduled.push(goroutine);
  if (!$schedulerLoopActive) {
    $schedulerLoopActive = true;
    setTimeout(function() {
      while (true) {
        var r = $scheduled.shift();
        if (r === undefined) {
          $schedulerLoopActive = false;
          break;
        }
        r();
      };
    }, 0);
  }
};

var $send = function(chan, value) {
  if (chan.$closed) {
    $throwRuntimeError("send on closed channel");
  }
  var queuedRecv = chan.$recvQueue.shift();
  if (queuedRecv !== undefined) {
    queuedRecv([value, true]);
    return;
  }
  if (chan.$buffer.length < chan.$capacity) {
    chan.$buffer.push(value);
    return;
  }

  var thisGoroutine = $curGoroutine;
  chan.$sendQueue.push(function() {
    $schedule(thisGoroutine);
    return value;
  });
  var blocked = false;
  var f = function() {
    if (blocked) {
      if (chan.$closed) {
        $throwRuntimeError("send on closed channel");
      }
      return;
    };
    blocked = true;
    $curGoroutine.asleep = true;
    throw null;
  };
  f.$blocking = true;
  return f;
};
var $recv = function(chan) {
  var queuedSend = chan.$sendQueue.shift();
  if (queuedSend !== undefined) {
    chan.$buffer.push(queuedSend());
  }
  var bufferedValue = chan.$buffer.shift();
  if (bufferedValue !== undefined) {
    return [bufferedValue, true];
  }
  if (chan.$closed) {
    return [chan.constructor.elem.zero(), false];
  }

  var thisGoroutine = $curGoroutine, value;
  var queueEntry = function(v) {
    value = v;
    $schedule(thisGoroutine);
  };
  chan.$recvQueue.push(queueEntry);
  var blocked = false;
  var f = function() {
    if (blocked) {
      return value;
    };
    blocked = true;
    $curGoroutine.asleep = true;
    throw null;
  };
  f.$blocking = true;
  return f;
};
var $close = function(chan) {
  if (chan.$closed) {
    $throwRuntimeError("close of closed channel");
  }
  chan.$closed = true;
  while (true) {
    var queuedSend = chan.$sendQueue.shift();
    if (queuedSend === undefined) {
      break;
    }
    queuedSend(); /* will panic because of closed channel */
  }
  while (true) {
    var queuedRecv = chan.$recvQueue.shift();
    if (queuedRecv === undefined) {
      break;
    }
    queuedRecv([chan.constructor.elem.zero(), false]);
  }
};
var $select = function(comms) {
  var ready = [], i;
  var selection = -1;
  for (i = 0; i < comms.length; i++) {
    var comm = comms[i];
    var chan = comm[0];
    switch (comm.length) {
    case 0: /* default */
      selection = i;
      break;
    case 1: /* recv */
      if (chan.$sendQueue.length !== 0 || chan.$buffer.length !== 0 || chan.$closed) {
        ready.push(i);
      }
      break;
    case 2: /* send */
      if (chan.$closed) {
        $throwRuntimeError("send on closed channel");
      }
      if (chan.$recvQueue.length !== 0 || chan.$buffer.length < chan.$capacity) {
        ready.push(i);
      }
      break;
    }
  }

  if (ready.length !== 0) {
    selection = ready[Math.floor(Math.random() * ready.length)];
  }
  if (selection !== -1) {
    var comm = comms[selection];
    switch (comm.length) {
    case 0: /* default */
      return [selection];
    case 1: /* recv */
      return [selection, $recv(comm[0])];
    case 2: /* send */
      $send(comm[0], comm[1]);
      return [selection];
    }
  }

  var entries = [];
  var thisGoroutine = $curGoroutine;
  var removeFromQueues = function() {
    for (i = 0; i < entries.length; i++) {
      var entry = entries[i];
      var queue = entry[0];
      var index = queue.indexOf(entry[1]);
      if (index !== -1) {
        queue.splice(index, 1);
      }
    }
  };
  for (i = 0; i < comms.length; i++) {
    (function(i) {
      var comm = comms[i];
      switch (comm.length) {
      case 1: /* recv */
        var queueEntry = function(value) {
          selection = [i, value];
          removeFromQueues();
          $schedule(thisGoroutine);
        };
        entries.push([comm[0].$recvQueue, queueEntry]);
        comm[0].$recvQueue.push(queueEntry);
        break;
      case 2: /* send */
        var queueEntry = function() {
          if (comm[0].$closed) {
            $throwRuntimeError("send on closed channel");
          }
          selection = [i];
          removeFromQueues();
          $schedule(thisGoroutine);
          return comm[1];
        };
        entries.push([comm[0].$sendQueue, queueEntry]);
        comm[0].$sendQueue.push(queueEntry);
        break;
      }
    })(i);
  }
  var blocked = false;
  var f = function() {
    if (blocked) {
      return selection;
    };
    blocked = true;
    $curGoroutine.asleep = true;
    throw null;
  };
  f.$blocking = true;
  return f;
};

var $needsExternalization = function(t) {
  switch (t.kind) {
    case $kindBool:
    case $kindInt:
    case $kindInt8:
    case $kindInt16:
    case $kindInt32:
    case $kindUint:
    case $kindUint8:
    case $kindUint16:
    case $kindUint32:
    case $kindUintptr:
    case $kindFloat32:
    case $kindFloat64:
      return false;
    case $kindInterface:
      return t !== $packages["github.com/gopherjs/gopherjs/js"].Object;
    default:
      return true;
  }
};

var $externalize = function(v, t) {
  switch (t.kind) {
  case $kindBool:
  case $kindInt:
  case $kindInt8:
  case $kindInt16:
  case $kindInt32:
  case $kindUint:
  case $kindUint8:
  case $kindUint16:
  case $kindUint32:
  case $kindUintptr:
  case $kindFloat32:
  case $kindFloat64:
    return v;
  case $kindInt64:
  case $kindUint64:
    return $flatten64(v);
  case $kindArray:
    if ($needsExternalization(t.elem)) {
      return $mapArray(v, function(e) { return $externalize(e, t.elem); });
    }
    return v;
  case $kindFunc:
    if (v === $throwNilPointerError) {
      return null;
    }
    if (v.$externalizeWrapper === undefined) {
      $checkForDeadlock = false;
      var convert = false;
      var i;
      for (i = 0; i < t.params.length; i++) {
        convert = convert || (t.params[i] !== $packages["github.com/gopherjs/gopherjs/js"].Object);
      }
      for (i = 0; i < t.results.length; i++) {
        convert = convert || $needsExternalization(t.results[i]);
      }
      v.$externalizeWrapper = v;
      if (convert) {
        v.$externalizeWrapper = function() {
          var args = [], i;
          for (i = 0; i < t.params.length; i++) {
            if (t.variadic && i === t.params.length - 1) {
              var vt = t.params[i].elem, varargs = [], j;
              for (j = i; j < arguments.length; j++) {
                varargs.push($internalize(arguments[j], vt));
              }
              args.push(new (t.params[i])(varargs));
              break;
            }
            args.push($internalize(arguments[i], t.params[i]));
          }
          var result = v.apply(this, args);
          switch (t.results.length) {
          case 0:
            return;
          case 1:
            return $externalize(result, t.results[0]);
          default:
            for (i = 0; i < t.results.length; i++) {
              result[i] = $externalize(result[i], t.results[i]);
            }
            return result;
          }
        };
      }
    }
    return v.$externalizeWrapper;
  case $kindInterface:
    if (v === $ifaceNil) {
      return null;
    }
    if (t === $packages["github.com/gopherjs/gopherjs/js"].Object || v.constructor.kind === undefined) {
      return v;
    }
    return $externalize(v.$val, v.constructor);
  case $kindMap:
    var m = {};
    var keys = $keys(v), i;
    for (i = 0; i < keys.length; i++) {
      var entry = v[keys[i]];
      m[$externalize(entry.k, t.key)] = $externalize(entry.v, t.elem);
    }
    return m;
  case $kindPtr:
    var o = {}, i;
    for (i = 0; i < t.methods.length; i++) {
      var m = t.methods[i];
      if (m[2] !== "") { /* not exported */
        continue;
      }
      (function(m) {
        o[m[1]] = $externalize(function() {
          return v[m[0]].apply(v, arguments);
        }, m[3]);
      })(m);
    }
    return o;
  case $kindSlice:
    if ($needsExternalization(t.elem)) {
      return $mapArray($sliceToArray(v), function(e) { return $externalize(e, t.elem); });
    }
    return $sliceToArray(v);
  case $kindString:
    if (v.search(/^[\x00-\x7F]*$/) !== -1) {
      return v;
    }
    var s = "", r, i;
    for (i = 0; i < v.length; i += r[1]) {
      r = $decodeRune(v, i);
      s += String.fromCharCode(r[0]);
    }
    return s;
  case $kindStruct:
    var timePkg = $packages["time"];
    if (timePkg && v.constructor === timePkg.Time.Ptr) {
      var milli = $div64(v.UnixNano(), new $Int64(0, 1000000));
      return new Date($flatten64(milli));
    }
    var o = {}, i;
    for (i = 0; i < t.fields.length; i++) {
      var f = t.fields[i];
      if (f[2] !== "") { /* not exported */
        continue;
      }
      o[f[1]] = $externalize(v[f[0]], f[3]);
    }
    return o;
  }
  $panic(new $String("cannot externalize " + t.string));
};

var $internalize = function(v, t, recv) {
  switch (t.kind) {
  case $kindBool:
    return !!v;
  case $kindInt:
    return parseInt(v);
  case $kindInt8:
    return parseInt(v) << 24 >> 24;
  case $kindInt16:
    return parseInt(v) << 16 >> 16;
  case $kindInt32:
    return parseInt(v) >> 0;
  case $kindUint:
    return parseInt(v);
  case $kindUint8:
    return parseInt(v) << 24 >>> 24;
  case $kindUint16:
    return parseInt(v) << 16 >>> 16;
  case $kindUint32:
  case $kindUintptr:
    return parseInt(v) >>> 0;
  case $kindInt64:
  case $kindUint64:
    return new t(0, v);
  case $kindFloat32:
  case $kindFloat64:
    return parseFloat(v);
  case $kindArray:
    if (v.length !== t.len) {
      $throwRuntimeError("got array with wrong size from JavaScript native");
    }
    return $mapArray(v, function(e) { return $internalize(e, t.elem); });
  case $kindFunc:
    return function() {
      var args = [], i;
      for (i = 0; i < t.params.length; i++) {
        if (t.variadic && i === t.params.length - 1) {
          var vt = t.params[i].elem, varargs = arguments[i], j;
          for (j = 0; j < varargs.$length; j++) {
            args.push($externalize(varargs.$array[varargs.$offset + j], vt));
          }
          break;
        }
        args.push($externalize(arguments[i], t.params[i]));
      }
      var result = v.apply(recv, args);
      switch (t.results.length) {
      case 0:
        return;
      case 1:
        return $internalize(result, t.results[0]);
      default:
        for (i = 0; i < t.results.length; i++) {
          result[i] = $internalize(result[i], t.results[i]);
        }
        return result;
      }
    };
  case $kindInterface:
    if (t === $packages["github.com/gopherjs/gopherjs/js"].Object) {
      return v;
    }
    if (v === null) {
      return $ifaceNil;
    }
    switch (v.constructor) {
    case Int8Array:
      return new ($sliceType($Int8))(v);
    case Int16Array:
      return new ($sliceType($Int16))(v);
    case Int32Array:
      return new ($sliceType($Int))(v);
    case Uint8Array:
      return new ($sliceType($Uint8))(v);
    case Uint16Array:
      return new ($sliceType($Uint16))(v);
    case Uint32Array:
      return new ($sliceType($Uint))(v);
    case Float32Array:
      return new ($sliceType($Float32))(v);
    case Float64Array:
      return new ($sliceType($Float64))(v);
    case Array:
      return $internalize(v, $sliceType($emptyInterface));
    case Boolean:
      return new $Bool(!!v);
    case Date:
      var timePkg = $packages["time"];
      if (timePkg) {
        return new timePkg.Time(timePkg.Unix(new $Int64(0, 0), new $Int64(0, v.getTime() * 1000000)));
      }
    case Function:
      var funcType = $funcType([$sliceType($emptyInterface)], [$packages["github.com/gopherjs/gopherjs/js"].Object], true);
      return new funcType($internalize(v, funcType));
    case Number:
      return new $Float64(parseFloat(v));
    case String:
      return new $String($internalize(v, $String));
    default:
      if ($global.Node && v instanceof $global.Node) {
        return v;
      }
      var mapType = $mapType($String, $emptyInterface);
      return new mapType($internalize(v, mapType));
    }
  case $kindMap:
    var m = new $Map();
    var keys = $keys(v), i;
    for (i = 0; i < keys.length; i++) {
      var key = $internalize(keys[i], t.key);
      m[key.$key ? key.$key() : key] = { k: key, v: $internalize(v[keys[i]], t.elem) };
    }
    return m;
  case $kindSlice:
    return new t($mapArray(v, function(e) { return $internalize(e, t.elem); }));
  case $kindString:
    v = String(v);
    if (v.search(/^[\x00-\x7F]*$/) !== -1) {
      return v;
    }
    var s = "", i;
    for (i = 0; i < v.length; i++) {
      s += $encodeRune(v.charCodeAt(i));
    }
    return s;
  default:
    $panic(new $String("cannot internalize " + t.string));
  }
};

$packages["github.com/gopherjs/gopherjs/js"] = (function() {
	var $pkg = {}, sliceType$2, ptrType, Object, Error, init;
	Object = $pkg.Object = $newType(8, $kindInterface, "js.Object", "Object", "github.com/gopherjs/gopherjs/js", null);
	Error = $pkg.Error = $newType(0, $kindStruct, "js.Error", "Error", "github.com/gopherjs/gopherjs/js", function(Object_) {
		this.$val = this;
		this.Object = Object_ !== undefined ? Object_ : null;
	});
	Error.Ptr.prototype.Error = function() {
		var err;
		err = this;
		return "JavaScript error: " + $internalize(err.Object.message, $String);
	};
	Error.prototype.Error = function() { return this.$val.Error(); };
	Error.Ptr.prototype.Stack = function() {
		var err;
		err = this;
		return $internalize(err.Object.stack, $String);
	};
	Error.prototype.Stack = function() { return this.$val.Stack(); };
	init = function() {
		var e;
		e = new Error.Ptr(null);
	};
	$pkg.$init = function() {
		$pkg.$init = function() {};
		/* */ var $r, $s = 0; var $init_js = function() { while (true) { switch ($s) { case 0:
		sliceType$2 = $sliceType($emptyInterface);
		ptrType = $ptrType(Error);
		Object.init([["Bool", "Bool", "", $funcType([], [$Bool], false)], ["Call", "Call", "", $funcType([$String, sliceType$2], [Object], true)], ["Delete", "Delete", "", $funcType([$String], [], false)], ["Float", "Float", "", $funcType([], [$Float64], false)], ["Get", "Get", "", $funcType([$String], [Object], false)], ["Index", "Index", "", $funcType([$Int], [Object], false)], ["Int", "Int", "", $funcType([], [$Int], false)], ["Int64", "Int64", "", $funcType([], [$Int64], false)], ["Interface", "Interface", "", $funcType([], [$emptyInterface], false)], ["Invoke", "Invoke", "", $funcType([sliceType$2], [Object], true)], ["Length", "Length", "", $funcType([], [$Int], false)], ["New", "New", "", $funcType([sliceType$2], [Object], true)], ["Set", "Set", "", $funcType([$String, $emptyInterface], [], false)], ["SetIndex", "SetIndex", "", $funcType([$Int, $emptyInterface], [], false)], ["Str", "Str", "", $funcType([], [$String], false)], ["Uint64", "Uint64", "", $funcType([], [$Uint64], false)], ["Unsafe", "Unsafe", "", $funcType([], [$Uintptr], false)]]);
		Error.methods = [["Bool", "Bool", "", $funcType([], [$Bool], false), 0], ["Call", "Call", "", $funcType([$String, sliceType$2], [Object], true), 0], ["Delete", "Delete", "", $funcType([$String], [], false), 0], ["Float", "Float", "", $funcType([], [$Float64], false), 0], ["Get", "Get", "", $funcType([$String], [Object], false), 0], ["Index", "Index", "", $funcType([$Int], [Object], false), 0], ["Int", "Int", "", $funcType([], [$Int], false), 0], ["Int64", "Int64", "", $funcType([], [$Int64], false), 0], ["Interface", "Interface", "", $funcType([], [$emptyInterface], false), 0], ["Invoke", "Invoke", "", $funcType([sliceType$2], [Object], true), 0], ["Length", "Length", "", $funcType([], [$Int], false), 0], ["New", "New", "", $funcType([sliceType$2], [Object], true), 0], ["Set", "Set", "", $funcType([$String, $emptyInterface], [], false), 0], ["SetIndex", "SetIndex", "", $funcType([$Int, $emptyInterface], [], false), 0], ["Str", "Str", "", $funcType([], [$String], false), 0], ["Uint64", "Uint64", "", $funcType([], [$Uint64], false), 0], ["Unsafe", "Unsafe", "", $funcType([], [$Uintptr], false), 0]];
		ptrType.methods = [["Bool", "Bool", "", $funcType([], [$Bool], false), 0], ["Call", "Call", "", $funcType([$String, sliceType$2], [Object], true), 0], ["Delete", "Delete", "", $funcType([$String], [], false), 0], ["Error", "Error", "", $funcType([], [$String], false), -1], ["Float", "Float", "", $funcType([], [$Float64], false), 0], ["Get", "Get", "", $funcType([$String], [Object], false), 0], ["Index", "Index", "", $funcType([$Int], [Object], false), 0], ["Int", "Int", "", $funcType([], [$Int], false), 0], ["Int64", "Int64", "", $funcType([], [$Int64], false), 0], ["Interface", "Interface", "", $funcType([], [$emptyInterface], false), 0], ["Invoke", "Invoke", "", $funcType([sliceType$2], [Object], true), 0], ["Length", "Length", "", $funcType([], [$Int], false), 0], ["New", "New", "", $funcType([sliceType$2], [Object], true), 0], ["Set", "Set", "", $funcType([$String, $emptyInterface], [], false), 0], ["SetIndex", "SetIndex", "", $funcType([$Int, $emptyInterface], [], false), 0], ["Stack", "Stack", "", $funcType([], [$String], false), -1], ["Str", "Str", "", $funcType([], [$String], false), 0], ["Uint64", "Uint64", "", $funcType([], [$Uint64], false), 0], ["Unsafe", "Unsafe", "", $funcType([], [$Uintptr], false), 0]];
		Error.init([["Object", "", "", Object, ""]]);
		init();
		/* */ } return; } }; $init_js.$blocking = true; return $init_js;
	};
	return $pkg;
})();
$packages["runtime"] = (function() {
	var $pkg = {}, js, ptrType$5, ptrType$6, ptrType$7, NotSupportedError, TypeAssertionError, errorString, init;
	js = $packages["github.com/gopherjs/gopherjs/js"];
	NotSupportedError = $pkg.NotSupportedError = $newType(0, $kindStruct, "runtime.NotSupportedError", "NotSupportedError", "runtime", function(Feature_) {
		this.$val = this;
		this.Feature = Feature_ !== undefined ? Feature_ : "";
	});
	TypeAssertionError = $pkg.TypeAssertionError = $newType(0, $kindStruct, "runtime.TypeAssertionError", "TypeAssertionError", "runtime", function(interfaceString_, concreteString_, assertedString_, missingMethod_) {
		this.$val = this;
		this.interfaceString = interfaceString_ !== undefined ? interfaceString_ : "";
		this.concreteString = concreteString_ !== undefined ? concreteString_ : "";
		this.assertedString = assertedString_ !== undefined ? assertedString_ : "";
		this.missingMethod = missingMethod_ !== undefined ? missingMethod_ : "";
	});
	errorString = $pkg.errorString = $newType(8, $kindString, "runtime.errorString", "errorString", "runtime", null);
	NotSupportedError.Ptr.prototype.Error = function() {
		var err;
		err = this;
		return "not supported by GopherJS: " + err.Feature;
	};
	NotSupportedError.prototype.Error = function() { return this.$val.Error(); };
	init = function() {
		var e;
		$throwRuntimeError = (function(msg) {
			$panic(new errorString(msg));
		});
		e = $ifaceNil;
		e = new TypeAssertionError.Ptr("", "", "", "");
		e = new NotSupportedError.Ptr("");
	};
	TypeAssertionError.Ptr.prototype.RuntimeError = function() {
	};
	TypeAssertionError.prototype.RuntimeError = function() { return this.$val.RuntimeError(); };
	TypeAssertionError.Ptr.prototype.Error = function() {
		var e, inter;
		e = this;
		inter = e.interfaceString;
		if (inter === "") {
			inter = "interface";
		}
		if (e.concreteString === "") {
			return "interface conversion: " + inter + " is nil, not " + e.assertedString;
		}
		if (e.missingMethod === "") {
			return "interface conversion: " + inter + " is " + e.concreteString + ", not " + e.assertedString;
		}
		return "interface conversion: " + e.concreteString + " is not " + e.assertedString + ": missing method " + e.missingMethod;
	};
	TypeAssertionError.prototype.Error = function() { return this.$val.Error(); };
	errorString.prototype.RuntimeError = function() {
		var e;
		e = this.$val;
	};
	$ptrType(errorString).prototype.RuntimeError = function() { return new errorString(this.$get()).RuntimeError(); };
	errorString.prototype.Error = function() {
		var e;
		e = this.$val;
		return "runtime error: " + e;
	};
	$ptrType(errorString).prototype.Error = function() { return new errorString(this.$get()).Error(); };
	$pkg.$init = function() {
		$pkg.$init = function() {};
		/* */ var $r, $s = 0; var $init_runtime = function() { while (true) { switch ($s) { case 0:
		$r = js.$init($BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
		ptrType$5 = $ptrType(NotSupportedError);
		ptrType$6 = $ptrType(TypeAssertionError);
		ptrType$7 = $ptrType(errorString);
		ptrType$5.methods = [["Error", "Error", "", $funcType([], [$String], false), -1]];
		NotSupportedError.init([["Feature", "Feature", "", $String, ""]]);
		ptrType$6.methods = [["Error", "Error", "", $funcType([], [$String], false), -1], ["RuntimeError", "RuntimeError", "", $funcType([], [], false), -1]];
		TypeAssertionError.init([["interfaceString", "interfaceString", "runtime", $String, ""], ["concreteString", "concreteString", "runtime", $String, ""], ["assertedString", "assertedString", "runtime", $String, ""], ["missingMethod", "missingMethod", "runtime", $String, ""]]);
		errorString.methods = [["Error", "Error", "", $funcType([], [$String], false), -1], ["RuntimeError", "RuntimeError", "", $funcType([], [], false), -1]];
		ptrType$7.methods = [["Error", "Error", "", $funcType([], [$String], false), -1], ["RuntimeError", "RuntimeError", "", $funcType([], [], false), -1]];
		init();
		/* */ } return; } }; $init_runtime.$blocking = true; return $init_runtime;
	};
	return $pkg;
})();
$packages["errors"] = (function() {
	var $pkg = {}, ptrType, errorString, New;
	errorString = $pkg.errorString = $newType(0, $kindStruct, "errors.errorString", "errorString", "errors", function(s_) {
		this.$val = this;
		this.s = s_ !== undefined ? s_ : "";
	});
	New = $pkg.New = function(text) {
		return new errorString.Ptr(text);
	};
	errorString.Ptr.prototype.Error = function() {
		var e;
		e = this;
		return e.s;
	};
	errorString.prototype.Error = function() { return this.$val.Error(); };
	$pkg.$init = function() {
		$pkg.$init = function() {};
		/* */ var $r, $s = 0; var $init_errors = function() { while (true) { switch ($s) { case 0:
		ptrType = $ptrType(errorString);
		ptrType.methods = [["Error", "Error", "", $funcType([], [$String], false), -1]];
		errorString.init([["s", "s", "errors", $String, ""]]);
		/* */ } return; } }; $init_errors.$blocking = true; return $init_errors;
	};
	return $pkg;
})();
$packages["sort"] = (function() {
	var $pkg = {}, Search;
	Search = $pkg.Search = function(n, f) {
		var _q, _tmp, _tmp$1, h, i, j;
		_tmp = 0; _tmp$1 = n; i = _tmp; j = _tmp$1;
		while (i < j) {
			h = i + (_q = ((j - i >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			if (!f(h)) {
				i = h + 1 >> 0;
			} else {
				j = h;
			}
		}
		return i;
	};
	$pkg.$init = function() {
		$pkg.$init = function() {};
		/* */ var $r, $s = 0; var $init_sort = function() { while (true) { switch ($s) { case 0:
		/* */ } return; } }; $init_sort.$blocking = true; return $init_sort;
	};
	return $pkg;
})();
$packages["ninchatclient"] = (function() {
	var $pkg = {}, errors, js, sort, sliceType, mapType, sliceType$1, ptrType, funcType, ptrType$1, chanType, funcType$1, funcType$2, funcType$3, ptrType$2, sliceType$2, funcType$4, chanType$1, ptrType$3, sliceType$3, funcType$5, funcType$6, funcType$7, funcType$8, funcType$9, ptrType$4, ptrType$5, funcType$10, funcType$11, sliceType$4, sliceType$5, chanType$2, ptrType$6, ptrType$7, Action, Backoff, Deferred, Transport, Session, Time, Duration, Timer, WebSocket, module, sessionEventAckWindow, xhrType, xhrRequestHeaderSupport, GetAddress, GetEndpointHosts, GetSessionEventCredentials, GetEventFrames, GetEventAndActionId, IsEventLastReply, GetEventError, Call, Jitter, JitterDuration, JitterUint64, jsError, jsInvoke, Atob, ParseDataURI, NewArray, NewArrayBuffer, NewUint8Array, NewObject, ParseJSON, StringifyJSON, Random, SetTimeout, ClearTimeout, Log, LongPollTransport, longPollTransfer, longPollPing, longPollClose, main, Defer, NewSession, Now, NewTimer, Sleep, NewWebSocket, StringifyFrame, WebSocketTransport, webSocketHandshake, webSocketSend, webSocketReceive, init, XHR, XHR_JSON;
	errors = $packages["errors"];
	js = $packages["github.com/gopherjs/gopherjs/js"];
	sort = $packages["sort"];
	Action = $pkg.Action = $newType(0, $kindStruct, "main.Action", "Action", "ninchatclient", function(Id_, Header_, Payload_, Deferred_, name_) {
		this.$val = this;
		this.Id = Id_ !== undefined ? Id_ : new $Uint64(0, 0);
		this.Header = Header_ !== undefined ? Header_ : null;
		this.Payload = Payload_ !== undefined ? Payload_ : null;
		this.Deferred = Deferred_ !== undefined ? Deferred_ : ptrType$2.nil;
		this.name = name_ !== undefined ? name_ : "";
	});
	Backoff = $pkg.Backoff = $newType(0, $kindStruct, "main.Backoff", "Backoff", "ninchatclient", function(lastSlot_) {
		this.$val = this;
		this.lastSlot = lastSlot_ !== undefined ? lastSlot_ : 0;
	});
	Deferred = $pkg.Deferred = $newType(0, $kindStruct, "main.Deferred", "Deferred", "ninchatclient", function(resolve_, reject_, notify_) {
		this.$val = this;
		this.resolve = resolve_ !== undefined ? resolve_ : sliceType$2.nil;
		this.reject = reject_ !== undefined ? reject_ : sliceType$2.nil;
		this.notify = notify_ !== undefined ? notify_ : sliceType$2.nil;
	});
	Transport = $pkg.Transport = $newType(4, $kindFunc, "main.Transport", "Transport", "ninchatclient", null);
	Session = $pkg.Session = $newType(0, $kindStruct, "main.Session", "Session", "ninchatclient", function(onSessionEvent_, onEvent_, onConnState_, onConnActive_, onLog_, address_, forceLongPoll_, sessionParams_, sessionId_, latestConnState_, latestConnActive_, lastActionId_, sendNotify_, sendBuffer_, numSent_, sendEventAck_, receivedEventId_, ackedEventId_, closeNotify_, closed_, stopped_) {
		this.$val = this;
		this.onSessionEvent = onSessionEvent_ !== undefined ? onSessionEvent_ : null;
		this.onEvent = onEvent_ !== undefined ? onEvent_ : null;
		this.onConnState = onConnState_ !== undefined ? onConnState_ : null;
		this.onConnActive = onConnActive_ !== undefined ? onConnActive_ : null;
		this.onLog = onLog_ !== undefined ? onLog_ : null;
		this.address = address_ !== undefined ? address_ : "";
		this.forceLongPoll = forceLongPoll_ !== undefined ? forceLongPoll_ : false;
		this.sessionParams = sessionParams_ !== undefined ? sessionParams_ : null;
		this.sessionId = sessionId_ !== undefined ? sessionId_ : null;
		this.latestConnState = latestConnState_ !== undefined ? latestConnState_ : "";
		this.latestConnActive = latestConnActive_ !== undefined ? latestConnActive_ : new Time(0, 0);
		this.lastActionId = lastActionId_ !== undefined ? lastActionId_ : new $Uint64(0, 0);
		this.sendNotify = sendNotify_ !== undefined ? sendNotify_ : chanType$1.nil;
		this.sendBuffer = sendBuffer_ !== undefined ? sendBuffer_ : sliceType$3.nil;
		this.numSent = numSent_ !== undefined ? numSent_ : 0;
		this.sendEventAck = sendEventAck_ !== undefined ? sendEventAck_ : false;
		this.receivedEventId = receivedEventId_ !== undefined ? receivedEventId_ : new $Uint64(0, 0);
		this.ackedEventId = ackedEventId_ !== undefined ? ackedEventId_ : new $Uint64(0, 0);
		this.closeNotify = closeNotify_ !== undefined ? closeNotify_ : chanType$1.nil;
		this.closed = closed_ !== undefined ? closed_ : false;
		this.stopped = stopped_ !== undefined ? stopped_ : false;
	});
	Time = $pkg.Time = $newType(8, $kindInt64, "main.Time", "Time", "ninchatclient", null);
	Duration = $pkg.Duration = $newType(8, $kindInt64, "main.Duration", "Duration", "ninchatclient", null);
	Timer = $pkg.Timer = $newType(0, $kindStruct, "main.Timer", "Timer", "ninchatclient", function(C_, id_) {
		this.$val = this;
		this.C = C_ !== undefined ? C_ : chanType$1.nil;
		this.id = id_ !== undefined ? id_ : null;
	});
	WebSocket = $pkg.WebSocket = $newType(0, $kindStruct, "main.WebSocket", "WebSocket", "ninchatclient", function(Notify_, impl_, open_, error_, buffer_) {
		this.$val = this;
		this.Notify = Notify_ !== undefined ? Notify_ : chanType$1.nil;
		this.impl = impl_ !== undefined ? impl_ : null;
		this.open = open_ !== undefined ? open_ : false;
		this.error = error_ !== undefined ? error_ : $ifaceNil;
		this.buffer = buffer_ !== undefined ? buffer_ : sliceType$2.nil;
	});
	Action.Ptr.prototype.Name = function() {
		var a;
		a = this;
		if (a.name === "") {
			a.name = $internalize(a.Header.action, $String);
		}
		return a.name;
	};
	Action.prototype.Name = function() { return this.$val.Name(); };
	GetAddress = $pkg.GetAddress = function(address) {
		if (address === undefined || address === null) {
			return "api.ninchat.com";
		} else {
			return $internalize(address, $String);
		}
	};
	GetEndpointHosts = $pkg.GetEndpointHosts = function(response) {
		var $deferred = [], $err = null, _tuple, endpoint, err = $ifaceNil, hosts = sliceType.nil, i, jsHosts;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			var e;
			e = jsError($recover());
			if (!($interfaceIsEqual(e, $ifaceNil))) {
				err = e;
			}
		}), []]);
		_tuple = ParseJSON(response); endpoint = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			return [hosts, err];
		}
		jsHosts = endpoint.hosts;
		if ($parseInt(jsHosts.length) === 0) {
			err = errors.New("endpoint hosts array is empty");
			return [hosts, err];
		}
		hosts = sliceType.make($parseInt(jsHosts.length));
		i = 0;
		while (i < $parseInt(jsHosts.length)) {
			(i < 0 || i >= hosts.$length) ? $throwRuntimeError("index out of range") : hosts.$array[hosts.$offset + i] = $internalize(jsHosts[i], $String);
			i = i + (1) >> 0;
		}
		return [hosts, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [hosts, err]; }
	};
	GetSessionEventCredentials = $pkg.GetSessionEventCredentials = function(header) {
		var $deferred = [], $err = null, err = $ifaceNil, eventId = new $Uint64(0, 0), object, ok = false, sessionId = null, userAuth = null, userId = null;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			var e;
			e = jsError($recover());
			if (!($interfaceIsEqual(e, $ifaceNil))) {
				err = e;
			}
		}), []]);
		if (!($internalize(header.event, $String) === "session_created")) {
			return [userId, userAuth, sessionId, eventId, ok, err];
		}
		userId = header.user_id;
		object = header.user_auth;
		if (!(object === undefined)) {
			userAuth = object;
		}
		sessionId = header.session_id;
		eventId = $internalize(header.event_id, $Uint64);
		ok = true;
		return [userId, userAuth, sessionId, eventId, ok, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [userId, userAuth, sessionId, eventId, ok, err]; }
	};
	GetEventFrames = $pkg.GetEventFrames = function(header) {
		var $deferred = [], $err = null, err = $ifaceNil, frames = 0, object;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		object = header.frames;
		if (!(object === undefined)) {
			frames = $parseInt(object) >> 0;
			if (frames < 0) {
				frames = 0;
			}
		}
		return [frames, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [frames, err]; }
	};
	GetEventAndActionId = $pkg.GetEventAndActionId = function(header) {
		var $deferred = [], $err = null, actionId = new $Uint64(0, 0), err = $ifaceNil, eventId = new $Uint64(0, 0), object, object$1;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		object = header.event_id;
		if (!(object === undefined)) {
			eventId = $internalize(object, $Uint64);
		}
		object$1 = header.action_id;
		if (!(object$1 === undefined)) {
			actionId = $internalize(object$1, $Uint64);
		}
		return [eventId, actionId, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [eventId, actionId, err]; }
	};
	IsEventLastReply = $pkg.IsEventLastReply = function(header, action) {
		var $deferred = [], $err = null, channels, err = $ifaceNil, historyLength, lastReply = false, users;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		lastReply = true;
		historyLength = header.history_length;
		if (!(historyLength === undefined)) {
			if (($parseInt(historyLength) >> 0) > 0) {
				lastReply = false;
			}
		}
		if (action.name === "search") {
			users = header.users;
			channels = header.channels;
			if (!(users === undefined) || !(channels === undefined)) {
				lastReply = false;
			}
		}
		return [lastReply, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [lastReply, err]; }
	};
	GetEventError = $pkg.GetEventError = function(header) {
		var $deferred = [], $err = null, _ref, err = $ifaceNil, errorReason = "", errorType = "", object, sessionLost = false;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			var e;
			e = jsError($recover());
			if (!($interfaceIsEqual(e, $ifaceNil))) {
				err = e;
			}
		}), []]);
		if (!($internalize(header.event, $String) === "error")) {
			return [errorType, errorReason, sessionLost, err];
		}
		errorType = $internalize(header.error_type, $String);
		object = header.error_reason;
		if (!(object === undefined)) {
			errorReason = $internalize(object, $String);
		}
		_ref = errorType;
		if (_ref === "session_not_found") {
			sessionLost = true;
			if (!(errorReason === "")) {
				err = errors.New("error: " + errorType + " (" + errorReason + ")");
			} else {
				err = errors.New("error: " + errorType);
			}
		} else if (_ref === "connection_superseded" || _ref === "message_has_too_many_parts" || _ref === "message_part_too_long" || _ref === "message_too_long" || _ref === "request_malformed") {
			if (!(errorReason === "")) {
				err = errors.New("error: " + errorType + " (" + errorReason + ")");
			} else {
				err = errors.New("error: " + errorType);
			}
		}
		return [errorType, errorReason, sessionLost, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [errorType, errorReason, sessionLost, err]; }
	};
	Backoff.Ptr.prototype.Success = function() {
		var b;
		b = this;
		b.lastSlot = 0;
	};
	Backoff.prototype.Success = function() { return this.$val.Success(); };
	Backoff.Ptr.prototype.Failure = function(maxDelay) {
		var b, delay = new Duration(0, 0);
		b = this;
		if (b.lastSlot > 0) {
			delay = new Duration(0, Jitter($flatten64(maxDelay) * b.lastSlot / 1024, -0.5));
		}
		if (b.lastSlot < 1023) {
			b.lastSlot = ((((b.lastSlot + 1 >> 0)) << 1 >> 0)) - 1 >> 0;
		}
		return delay;
	};
	Backoff.prototype.Failure = function(maxDelay) { return this.$val.Failure(maxDelay); };
	Call = $pkg.Call = function(header, onLog, address) {
		var _tuple, deferred, promise = $ifaceNil, url;
		url = "https://" + GetAddress(address) + "/v2/call";
		_tuple = Defer(); deferred = _tuple[0]; promise = new mapType(_tuple[1]);
		$go((function($b) {
			var $args = arguments, $r, $s = 0, $this = this, _r, _tuple$1, _tuple$2, _tuple$3, channel, err, event, events, ok, response;
			/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
			_tuple$1 = XHR_JSON(url, header, JitterDuration(new Duration(0, 11000), 0.1)); channel = _tuple$1[0]; err = _tuple$1[1];
			if (!($interfaceIsEqual(err, $ifaceNil))) {
				Log("NinchatClient.call onLog callback", onLog, new sliceType$1([new $String("call:"), err]));
				deferred.Reject(new sliceType$1([]));
				return;
			}
			_r = $recv(channel, $BLOCKING); /* */ $s = 1; case 1: if (_r && _r.$blocking) { _r = _r(); }
			_tuple$2 = _r; response = _tuple$2[0]; ok = _tuple$2[1];
			if (response === "") {
				if (ok) {
					Log("NinchatClient.call onLog callback", onLog, new sliceType$1([new $String("call error")]));
				} else {
					Log("NinchatClient.call onLog callback", onLog, new sliceType$1([new $String("call timeout")]));
				}
				deferred.Reject(new sliceType$1([]));
				return;
			}
			_tuple$3 = ParseJSON(response); event = _tuple$3[0]; err = _tuple$3[1];
			if (!($interfaceIsEqual(err, $ifaceNil))) {
				Log("NinchatClient.call onLog callback", onLog, new sliceType$1([new $String("call response:"), err]));
				deferred.Reject(new sliceType$1([]));
				return;
			}
			events = NewArray();
			events.push(event);
			deferred.Resolve(new sliceType$1([events]));
			/* */ case -1: } return; } }; $f.$blocking = true; return $f;
		}), []);
		return promise;
	};
	Jitter = $pkg.Jitter = function(x, scale) {
		return x + x * scale * Random();
	};
	JitterDuration = $pkg.JitterDuration = function(d, scale) {
		return new Duration(0, Jitter($flatten64(d), scale));
	};
	JitterUint64 = $pkg.JitterUint64 = function(n, scale) {
		return new $Uint64(0, Jitter($flatten64(n), scale));
	};
	jsError = function(x) {
		var _tuple, err = $ifaceNil, jsErr, msg, ok;
		if ($interfaceIsEqual(x, $ifaceNil)) {
			return err;
		}
		_tuple = $assertType(x, ptrType, true); jsErr = _tuple[0]; ok = _tuple[1];
		if (ok) {
			msg = $internalize(jsErr.Object.message, $String);
			if (msg === "") {
				msg = "error";
			}
			err = errors.New(msg);
			return err;
		}
		err = $assertType(x, $error);
		return err;
	};
	jsInvoke = function(name, function$1, args) {
		var $deferred = [], $err = null, ok = false;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			var err;
			err = jsError($recover());
			if (!($interfaceIsEqual(err, $ifaceNil))) {
				console.log(name + " invocation error: " + err.Error());
			}
		}), []]);
		function$1.apply(undefined, $externalize(args, sliceType$1));
		ok = true;
		return ok;
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return ok; }
	};
	Atob = $pkg.Atob = function(string) {
		var $deferred = [], $err = null, binary = null, err = $ifaceNil;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		binary = $global.atob(string);
		return [binary, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [binary, err]; }
	};
	ParseDataURI = $pkg.ParseDataURI = function(string) {
		var $deferred = [], $err = null, base64 = null, err = $ifaceNil;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		base64 = string.split($externalize(",", $String))[1];
		return [base64, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [base64, err]; }
	};
	NewArray = $pkg.NewArray = function() {
		return new ($global.Array)();
	};
	NewArrayBuffer = $pkg.NewArrayBuffer = function(length) {
		return new ($global.ArrayBuffer)(length);
	};
	NewUint8Array = $pkg.NewUint8Array = function(arrayBuffer) {
		return new ($global.Uint8Array)(arrayBuffer);
	};
	NewObject = $pkg.NewObject = function() {
		return new ($global.Object)();
	};
	ParseJSON = $pkg.ParseJSON = function(json) {
		var $deferred = [], $err = null, err = $ifaceNil, object = null;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		object = $global.JSON.parse($externalize(json, $String));
		return [object, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [object, err]; }
	};
	StringifyJSON = $pkg.StringifyJSON = function(object) {
		var $deferred = [], $err = null, err = $ifaceNil, json = "";
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		json = $internalize($global.JSON.stringify($externalize(object, $emptyInterface)), $String);
		return [json, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [json, err]; }
	};
	Random = $pkg.Random = function() {
		return $parseFloat($global.Math.random());
	};
	SetTimeout = $pkg.SetTimeout = function(callback, timeout) {
		var id = null;
		id = $global.setTimeout($externalize(callback, funcType), $externalize(timeout, Duration));
		return id;
	};
	ClearTimeout = $pkg.ClearTimeout = function(id) {
		$global.clearTimeout(id);
	};
	Log = $pkg.Log = function(logInvocationName, onLog, tokens) {
		var _i, _ref, _tuple, _tuple$1, message, ok, ok$1, str, x, y, y$1;
		if (onLog === undefined || onLog === null) {
			return;
		}
		message = "";
		_ref = tokens;
		_i = 0;
		while (_i < _ref.$length) {
			x = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			str = "?";
			_tuple = $assertType(x, $String, true); y = _tuple[0]; ok = _tuple[1];
			if (ok) {
				str = y;
			} else {
				_tuple$1 = $assertType(x, $error, true); y$1 = _tuple$1[0]; ok$1 = _tuple$1[1];
				if (ok$1) {
					str = y$1.Error();
				}
			}
			if (message.length > 0) {
				message = message + (" ");
			}
			message = message + (str);
			_i++;
		}
		while (message.length > 0 && (message.charCodeAt((message.length - 1 >> 0)) === 32)) {
			message = message.substring(0, (message.length - 1 >> 0));
		}
		jsInvoke(logInvocationName, onLog, new sliceType$1([new $String(message)]));
	};
	LongPollTransport = $pkg.LongPollTransport = function(s, host, $b) {
		var $args = arguments, $deferred = [], $err = null, $r, $s = 0, $this = this, _r, _selection, _tuple, _tuple$1, _tuple$2, array, connWorked = false, creator, err, err$1, err$2, gotOnline = false, header, header$1, ok, response, url;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_LongPollTransport = function() { try { $deferFrames.push($deferred); s: while (true) { switch ($s) { case 0:
		$deferred.push([(function() {
			var err;
			err = jsError($recover());
			if (!($interfaceIsEqual(err, $ifaceNil))) {
				s.log(new sliceType$1([new $String("poll:"), err]));
			}
		}), [$BLOCKING]]);
		url = "https://" + host + "/v2/poll";
		/* if (s.sessionId === null) { */ if (s.sessionId === null) {} else { $s = 1; continue; }
			s.log(new sliceType$1([new $String("session creation")]));
			header = s.makeCreateSessionAction();
			_tuple = XHR_JSON(url, header, JitterDuration(new Duration(0, 13000), 0.2)); creator = _tuple[0]; err = _tuple[1];
			if (!($interfaceIsEqual(err, $ifaceNil))) {
				s.log(new sliceType$1([new $String("session creation:"), err]));
				return [connWorked, gotOnline];
			}
			_r = $select([[creator], [s.closeNotify]], $BLOCKING); /* */ $s = 3; case 3: if (_r && _r.$blocking) { _r = _r(); }
			_selection = _r;
			if (_selection[0] === 0) {
				_tuple$1 = _selection[1]; response = _tuple$1[0]; ok = _tuple$1[1];
				if (!ok) {
					s.log(new sliceType$1([new $String("session creation timeout")]));
					return [connWorked, gotOnline];
				} else if (response === "") {
					s.log(new sliceType$1([new $String("session creation error")]));
					return [connWorked, gotOnline];
				}
				_tuple$2 = ParseJSON(response); array = _tuple$2[0]; err$1 = _tuple$2[1];
				if (!($interfaceIsEqual(err$1, $ifaceNil))) {
					s.log(new sliceType$1([new $String("session creation response:"), err$1]));
					return [connWorked, gotOnline];
				}
				header$1 = array[0];
				if (!s.handleSessionEvent(header$1)) {
					return [connWorked, gotOnline];
				}
			} else if (_selection[0] === 1) {
				longPollClose(s, url);
				return [connWorked, gotOnline];
			}
			connWorked = true;
			gotOnline = true;
			s.connState("connected");
			s.connActive();
		/* } else { */ $s = 2; continue; case 1: 
			s.log(new sliceType$1([new $String("session resumption")]));
			err$2 = longPollPing(s, url);
			if (!($interfaceIsEqual(err$2, $ifaceNil))) {
				s.log(new sliceType$1([new $String("session resumption:"), err$2]));
				return [connWorked, gotOnline];
			}
		/* } */ case 2:
		$r = longPollTransfer(s, url, new ptrType$1(function() { return connWorked; }, function($v) { connWorked = $v; }), new ptrType$1(function() { return gotOnline; }, function($v) { gotOnline = $v; }), $BLOCKING); /* */ $s = 4; case 4: if ($r && $r.$blocking) { $r = $r(); }
		return [connWorked, gotOnline];
		/* */ case -1: } return; } } catch(err) { $err = err; } finally { $deferFrames.pop(); if ($curGoroutine.asleep && !$jumpToDefer) { throw null; } $s = -1; $callDeferred($deferred, $err); return [connWorked, gotOnline]; } }; $blocking_LongPollTransport.$blocking = true; return $blocking_LongPollTransport;
	};
	longPollTransfer = function(s, url, connWorked, gotOnline, $b) {
		var $args = arguments, $r, $s = 0, $this = this, _r, _selection, _tuple, _tuple$1, _tuple$2, _tuple$3, _tuple$4, _tuple$5, _tuple$6, _tuple$7, _tuple$8, _tuple$9, ackedActionId, action, array, base64, channel, err, err$1, err$2, err$3, err$4, err$5, failures, frame, header, header$1, i, json, object, ok, ok$1, payload, payload$1, poller, request, response, sender, sending, sendingId, sessionLost, x, x$1, x$2;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_longPollTransfer = function() { s: while (true) { switch ($s) { case 0:
		poller = chanType.nil;
		sender = chanType.nil;
		sendingId = new $Uint64(0, 0);
		failures = 0;
		s.numSent = 0;
		/* while (failures < 2) { */ case 1: if(!(failures < 2)) { $s = 2; continue; }
			if (poller === chanType.nil) {
				err = $ifaceNil;
				header = s.makeResumeSessionAction(true);
				_tuple = XHR_JSON(url, header, JitterDuration(new Duration(0, 64000), 0.2)); poller = _tuple[0]; err = _tuple[1];
				if (!($interfaceIsEqual(err, $ifaceNil))) {
					s.log(new sliceType$1([new $String("poll:"), err]));
					return;
				}
			}
			if (sender === chanType.nil && s.numSent < s.sendBuffer.$length) {
				action = (x = s.sendBuffer, x$1 = s.numSent, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1]));
				if (!(action.Payload === null)) {
					payload = null;
					err$1 = $ifaceNil;
					frame = action.Payload[0];
					if ($internalize(action.Header.action, $String) === "update_user") {
						_tuple$1 = ParseDataURI(frame); base64 = _tuple$1[0]; err$2 = _tuple$1[1];
						if (!($interfaceIsEqual(err$2, $ifaceNil))) {
							s.log(new sliceType$1([new $String("send:"), err$2]));
							return;
						}
						payload = NewArray();
						payload.push(base64);
					} else {
						_tuple$2 = ParseJSON($internalize(frame, $String)); payload = _tuple$2[0]; err$1 = _tuple$2[1];
						if (!($interfaceIsEqual(err$1, $ifaceNil))) {
							s.log(new sliceType$1([new $String("send:"), err$1]));
							return;
						}
					}
					action.Header.payload = payload;
				}
				action.Header.session_id = s.sessionId;
				_tuple$3 = StringifyJSON(action.Header); request = _tuple$3[0]; err$3 = _tuple$3[1];
				delete action.Header[$externalize("session_id", $String)];
				delete action.Header[$externalize("payload", $String)];
				if (!($interfaceIsEqual(err$3, $ifaceNil))) {
					s.log(new sliceType$1([new $String("send:"), err$3]));
					return;
				}
				_tuple$4 = XHR(url, request, JitterDuration(new Duration(0, 7000), 0.2)); channel = _tuple$4[0]; err$3 = _tuple$4[1];
				if (!($interfaceIsEqual(err$3, $ifaceNil))) {
					s.log(new sliceType$1([new $String("send:"), err$3]));
					return;
				}
				if ((x$2 = action.Id, (x$2.$high === 0 && x$2.$low === 0))) {
					s.sendBuffer = $appendSlice($subslice(s.sendBuffer, 0, s.numSent), $subslice(s.sendBuffer, (s.numSent + 1 >> 0)));
				} else {
					sender = channel;
					sendingId = action.Id;
				}
			}
			response = "";
			ok = false;
			_r = $select([[poller], [sender], [s.sendNotify], [s.closeNotify]], $BLOCKING); /* */ $s = 3; case 3: if (_r && _r.$blocking) { _r = _r(); }
			_selection = _r;
			if (_selection[0] === 0) {
				_tuple$5 = _selection[1]; response = _tuple$5[0]; ok = _tuple$5[1];
				if (!ok) {
					s.log(new sliceType$1([new $String("poll timeout")]));
				} else if (response === "") {
					s.log(new sliceType$1([new $String("poll error")]));
				}
				poller = chanType.nil;
				s.connActive();
			} else if (_selection[0] === 1) {
				_tuple$6 = _selection[1]; response = _tuple$6[0]; ok = _tuple$6[1];
				if (!ok) {
					s.log(new sliceType$1([new $String("send timeout")]));
				} else if (response === "") {
					s.log(new sliceType$1([new $String("send error")]));
				} else if ((sendingId.$high > 0 || (sendingId.$high === 0 && sendingId.$low > 0))) {
					s.numSent = s.numSent + (1) >> 0;
				}
				sender = chanType.nil;
				sendingId = new $Uint64(0, 0);
			} else if (_selection[0] === 2) {
				sending = _selection[1][0];
				if (!sending) {
					longPollClose(s, url);
					return;
				}
				/* continue; */ $s = 1; continue;
			} else if (_selection[0] === 3) {
				longPollClose(s, url);
				return;
			}
			array = null;
			if (!(response === "")) {
				err$4 = $ifaceNil;
				_tuple$7 = ParseJSON(response); array = _tuple$7[0]; err$4 = _tuple$7[1];
				if (!($interfaceIsEqual(err$4, $ifaceNil))) {
					s.log(new sliceType$1([new $String("response:"), err$4]));
				}
			}
			if (array === null) {
				failures = failures + (1) >> 0;
				s.numSent = 0;
				/* continue; */ $s = 1; continue;
			}
			failures = 0;
			connWorked.$set(true);
			i = 0;
			while (i < $parseInt(array.length)) {
				header$1 = array[i];
				payload$1 = NewArray();
				object = header$1.payload;
				if (!(object === undefined)) {
					_tuple$8 = StringifyJSON(object); json = _tuple$8[0]; err$5 = _tuple$8[1];
					if (!($interfaceIsEqual(err$5, $ifaceNil))) {
						s.log(new sliceType$1([new $String("poll payload:"), err$5]));
						return;
					}
					payload$1.push($externalize(json, $String));
				}
				_tuple$9 = s.handleEvent(header$1, payload$1); ackedActionId = _tuple$9[0]; sessionLost = _tuple$9[1]; ok$1 = _tuple$9[3];
				if ((sendingId.$high > 0 || (sendingId.$high === 0 && sendingId.$low > 0)) && (sendingId.$high < ackedActionId.$high || (sendingId.$high === ackedActionId.$high && sendingId.$low <= ackedActionId.$low))) {
					sendingId = new $Uint64(0, 0);
					s.numSent = s.numSent + (1) >> 0;
				}
				if (!ok$1) {
					if (sessionLost) {
						gotOnline.$set(true);
					}
					return;
				}
				if (!gotOnline.$get()) {
					gotOnline.$set(true);
					s.connState("connected");
				}
				i = i + (1) >> 0;
			}
		/* } */ $s = 1; continue; case 2:
		return;
		/* */ case -1: } return; } }; $blocking_longPollTransfer.$blocking = true; return $blocking_longPollTransfer;
	};
	longPollPing = function(s, url) {
		var _key, _map, _tuple, err = $ifaceNil, header;
		header = (_map = new $Map(), _key = "action", _map[_key] = { k: _key, v: new $String("ping") }, _key = "session_id", _map[_key] = { k: _key, v: s.sessionId }, _map);
		_tuple = XHR_JSON(url, new mapType(header), JitterDuration(new Duration(0, 7000), 0.9)); err = _tuple[1];
		return err;
	};
	longPollClose = function(s, url) {
		var _key, _map, _tuple, err, header;
		header = (_map = new $Map(), _key = "action", _map[_key] = { k: _key, v: new $String("close_session") }, _key = "session_id", _map[_key] = { k: _key, v: s.sessionId }, _map);
		_tuple = XHR_JSON(url, new mapType(header), JitterDuration(new Duration(0, 7000), 0.9)); err = _tuple[1];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			s.log(new sliceType$1([new $String("send:"), err]));
		}
	};
	main = function() {
		module.call = $externalize(Call, funcType$1);
		module.newSession = $externalize(NewSession, funcType$2);
		module.stringifyFrame = $externalize(StringifyFrame, funcType$3);
		$global.NinchatClient = module;
	};
	Defer = $pkg.Defer = function() {
		var _key, _map, d = ptrType$2.nil, promise = false;
		d = new Deferred.Ptr(sliceType$2.nil, sliceType$2.nil, sliceType$2.nil);
		promise = (_map = new $Map(), _key = "then", _map[_key] = { k: _key, v: new funcType$4($methodVal(d, "then")) }, _map);
		return [d, promise];
	};
	Deferred.Ptr.prototype.then = function(resolve, reject, notify) {
		var d;
		d = this;
		if (!(resolve === undefined) && !(resolve === null)) {
			d.resolve = $append(d.resolve, resolve);
		}
		if (!(reject === undefined) && !(reject === null)) {
			d.reject = $append(d.reject, reject);
		}
		if (!(notify === undefined) && !(notify === null)) {
			d.notify = $append(d.notify, notify);
		}
	};
	Deferred.prototype.then = function(resolve, reject, notify) { return this.$val.then(resolve, reject, notify); };
	Deferred.Ptr.prototype.Resolve = function(args) {
		var _i, _ref, callback, d;
		d = this;
		_ref = d.resolve;
		_i = 0;
		while (_i < _ref.$length) {
			callback = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			jsInvoke("NinchatClient promise resolve callback", callback, args);
			_i++;
		}
	};
	Deferred.prototype.Resolve = function(args) { return this.$val.Resolve(args); };
	Deferred.Ptr.prototype.Reject = function(args) {
		var _i, _ref, callback, d;
		d = this;
		_ref = d.reject;
		_i = 0;
		while (_i < _ref.$length) {
			callback = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			jsInvoke("NinchatClient promise reject callback", callback, args);
			_i++;
		}
	};
	Deferred.prototype.Reject = function(args) { return this.$val.Reject(args); };
	Deferred.Ptr.prototype.Notify = function(args) {
		var _i, _ref, callback, d;
		d = this;
		_ref = d.notify;
		_i = 0;
		while (_i < _ref.$length) {
			callback = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			jsInvoke("NinchatClient promise notify callback", callback, args);
			_i++;
		}
	};
	Deferred.prototype.Notify = function(args) { return this.$val.Notify(args); };
	NewSession = $pkg.NewSession = function() {
		var _key, _map, s;
		s = new Session.Ptr(null, null, null, null, null, "api.ninchat.com", false, null, null, "", new Time(0, 0), new $Uint64(0, 0), chanType$1.nil, sliceType$3.nil, 0, false, new $Uint64(0, 0), new $Uint64(0, 0), chanType$1.nil, false, true);
		return (_map = new $Map(), _key = "onSessionEvent", _map[_key] = { k: _key, v: new funcType$5($methodVal(s, "OnSessionEvent")) }, _key = "onEvent", _map[_key] = { k: _key, v: new funcType$5($methodVal(s, "OnEvent")) }, _key = "onConnState", _map[_key] = { k: _key, v: new funcType$5($methodVal(s, "OnConnState")) }, _key = "onConnActive", _map[_key] = { k: _key, v: new funcType$5($methodVal(s, "OnConnActive")) }, _key = "onLog", _map[_key] = { k: _key, v: new funcType$5($methodVal(s, "OnLog")) }, _key = "setParams", _map[_key] = { k: _key, v: new funcType$6($methodVal(s, "SetParams")) }, _key = "setTransport", _map[_key] = { k: _key, v: new funcType$7($methodVal(s, "SetTransport")) }, _key = "setAddress", _map[_key] = { k: _key, v: new funcType$8($methodVal(s, "SetAddress")) }, _key = "open", _map[_key] = { k: _key, v: new funcType($methodVal(s, "Open")) }, _key = "close", _map[_key] = { k: _key, v: new funcType($methodVal(s, "Close")) }, _key = "send", _map[_key] = { k: _key, v: new funcType$9($methodVal(s, "Send")) }, _map);
	};
	Session.Ptr.prototype.OnSessionEvent = function(callback) {
		var s;
		s = this;
		s.onSessionEvent = callback;
	};
	Session.prototype.OnSessionEvent = function(callback) { return this.$val.OnSessionEvent(callback); };
	Session.Ptr.prototype.OnEvent = function(callback) {
		var s;
		s = this;
		s.onEvent = callback;
	};
	Session.prototype.OnEvent = function(callback) { return this.$val.OnEvent(callback); };
	Session.Ptr.prototype.OnConnState = function(callback) {
		var s;
		s = this;
		if (callback === null) {
			callback = null;
		}
		s.onConnState = callback;
		if (!(s.onConnState === null) && !(s.latestConnState === "")) {
			jsInvoke("NinchatClient.Session onConnState callback", s.onConnState, new sliceType$1([new $String(s.latestConnState)]));
		}
	};
	Session.prototype.OnConnState = function(callback) { return this.$val.OnConnState(callback); };
	Session.Ptr.prototype.OnConnActive = function(callback) {
		var s, x;
		s = this;
		if (callback === null) {
			callback = null;
		}
		s.onConnActive = callback;
		if (!(s.onConnActive === null) && (x = s.latestConnActive, (x.$high > 0 || (x.$high === 0 && x.$low > 0)))) {
			jsInvoke("NinchatClient.Session onConnActive callback", s.onConnActive, new sliceType$1([s.latestConnActive]));
		}
	};
	Session.prototype.OnConnActive = function(callback) { return this.$val.OnConnActive(callback); };
	Session.Ptr.prototype.OnLog = function(callback) {
		var s;
		s = this;
		if (callback === null) {
			callback = null;
		}
		s.onLog = callback;
	};
	Session.prototype.OnLog = function(callback) { return this.$val.OnLog(callback); };
	Session.Ptr.prototype.SetParams = function(params) {
		var s, sessionId;
		s = this;
		if (params.message_types === undefined) {
			$panic(new $String("message_types parameter not defined"));
		}
		sessionId = params.session_id;
		if (!(sessionId === undefined) && !(sessionId === null)) {
			s.sessionId = sessionId;
		}
		delete params[$externalize("session_id", $String)];
		s.sessionParams = params;
		if (!(s.sendNotify === chanType$1.nil) && s.stopped) {
			$go($methodVal(s, "discover"), []);
		}
	};
	Session.prototype.SetParams = function(params) { return this.$val.SetParams(params); };
	Session.Ptr.prototype.SetTransport = function(name) {
		var _ref, s, string;
		s = this;
		if (name === null) {
			s.forceLongPoll = false;
			return;
		}
		string = $internalize(name, $String);
		_ref = string;
		if (_ref === "websocket") {
			$panic(new $String("websocket transport cannot be forced"));
		} else if (_ref === "longpoll") {
			s.forceLongPoll = true;
		} else {
			$panic(new $String("unknown transport: " + string));
		}
	};
	Session.prototype.SetTransport = function(name) { return this.$val.SetTransport(name); };
	Session.Ptr.prototype.SetAddress = function(address) {
		var s;
		s = this;
		s.address = GetAddress(address);
	};
	Session.prototype.SetAddress = function(address) { return this.$val.SetAddress(address); };
	Session.Ptr.prototype.Open = function() {
		var s;
		s = this;
		if (s.closed) {
			$panic(new $String("session already closed"));
		}
		if (!(s.sendNotify === chanType$1.nil)) {
			$panic(new $String("session already initialized"));
		}
		if (s.onSessionEvent === null) {
			$panic(new $String("onSessionEvent callback not defined"));
		}
		if (s.onEvent === null) {
			$panic(new $String("onEvent callback not defined"));
		}
		if (s.sessionParams === null) {
			$panic(new $String("session parameters not defined"));
		}
		s.sendNotify = new chanType$1(1);
		s.closeNotify = new chanType$1(1);
		s.stopped = false;
		$go($methodVal(s, "discover"), []);
	};
	Session.prototype.Open = function() { return this.$val.Open(); };
	Session.Ptr.prototype.Close = function() {
		var _i, _ref, action, s;
		s = this;
		if (s.closed) {
			return;
		}
		_ref = s.sendBuffer;
		_i = 0;
		while (_i < _ref.$length) {
			action = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (!(action.Deferred === ptrType$2.nil)) {
				action.Deferred.Reject(new sliceType$1([]));
			}
			_i++;
		}
		s.sendBuffer = sliceType$3.nil;
		s.numSent = 0;
		s.closed = true;
		s.stopped = true;
		$go((function($b) {
			var $args = arguments, $r, $s = 0, $this = this;
			/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
			$r = $send(s.closeNotify, true, $BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
			$close(s.sendNotify);
			/* */ case -1: } return; } }; $f.$blocking = true; return $f;
		}), []);
	};
	Session.prototype.Close = function() { return this.$val.Close(); };
	Session.Ptr.prototype.Send = function(header, payload) {
		var _tuple, action, promise = $ifaceNil, s, x, x$1;
		s = this;
		if (s.sendNotify === chanType$1.nil) {
			$panic(new $String("session not initialized"));
		}
		if (s.closed) {
			$panic(new $String("session already closed"));
		}
		if (payload === undefined || payload === null || ($parseInt(payload.length) === 0)) {
			payload = null;
		}
		action = new Action.Ptr(new $Uint64(0, 0), header, payload, ptrType$2.nil, "");
		if (header.action_id === null) {
			delete header[$externalize("action_id", $String)];
		} else {
			s.lastActionId = (x = s.lastActionId, x$1 = new $Uint64(0, 1), new $Uint64(x.$high + x$1.$high, x.$low + x$1.$low));
			action.Id = s.lastActionId;
			header.action_id = $externalize(action.Id, $Uint64);
			_tuple = Defer(); action.Deferred = _tuple[0]; promise = new mapType(_tuple[1]);
		}
		s.send(action);
		return promise;
	};
	Session.prototype.Send = function(header, payload) { return this.$val.Send(header, payload); };
	Session.Ptr.prototype.send = function(action) {
		var s;
		s = this;
		s.sendBuffer = $append(s.sendBuffer, action);
		$go((function() {
			var _selection;
			_selection = $select([[s.sendNotify, true], []]);
		}), []);
		return;
	};
	Session.prototype.send = function(action) { return this.$val.send(action); };
	Session.Ptr.prototype.sendAck = function() {
		var s;
		s = this;
		s.sendEventAck = true;
		$go((function() {
			var _selection;
			_selection = $select([[s.sendNotify, true], []]);
		}), []);
	};
	Session.prototype.sendAck = function() { return this.$val.sendAck(); };
	Session.Ptr.prototype.discover = function($b) {
		var $args = arguments, $deferred = [], $err = null, $r, $s = 0, $this = this, _r, _r$1, _r$2, _selection, _tuple, _tuple$1, _tuple$2, backoff, channel, delay, err, err$1, hosts, ok, response, s, url;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_discover = function() { try { $deferFrames.push($deferred); s: while (true) { switch ($s) { case 0:
		s = $this;
		s.log(new sliceType$1([new $String("opening")]));
		$deferred.push([$methodVal(s, "log"), [new sliceType$1([new $String("closed")]), $BLOCKING]]);
		$deferred.push([$methodVal(s, "connState"), ["disconnected", $BLOCKING]]);
		backoff = $clone(new Backoff.Ptr(), Backoff);
		/* while (!s.stopped) { */ case 1: if(!(!s.stopped)) { $s = 2; continue; }
			s.log(new sliceType$1([new $String("endpoint discovery")]));
			s.connState("connecting");
			url = "https://" + s.address + "/v2/endpoint";
			_tuple = XHR(url, "", JitterDuration(new Duration(0, 7000), 0.1)); channel = _tuple[0]; err = _tuple[1];
			/* if (!($interfaceIsEqual(err, $ifaceNil))) { */ if (!($interfaceIsEqual(err, $ifaceNil))) {} else { $s = 3; continue; }
				s.log(new sliceType$1([new $String("endpoint discovery:"), err]));
			/* } else { */ $s = 4; continue; case 3: 
				_r = $select([[channel], [s.closeNotify]], $BLOCKING); /* */ $s = 5; case 5: if (_r && _r.$blocking) { _r = _r(); }
				_selection = _r;
				/* if (_selection[0] === 0) { */ if (_selection[0] === 0) {} else if (_selection[0] === 1) { $s = 6; continue; } else { $s = 7; continue; }
					_tuple$1 = _selection[1]; response = _tuple$1[0]; ok = _tuple$1[1];
					/* if (!ok) { */ if (!ok) {} else if (response === "") { $s = 8; continue; } else { $s = 9; continue; }
						s.log(new sliceType$1([new $String("endpoint discovery timeout")]));
					/* } else if (response === "") { */ $s = 10; continue; case 8: 
						s.log(new sliceType$1([new $String("endpoint discovery error")]));
					/* } else { */ $s = 10; continue; case 9: 
						_tuple$2 = GetEndpointHosts(response); hosts = _tuple$2[0]; err$1 = _tuple$2[1];
						/* if (!($interfaceIsEqual(err$1, $ifaceNil))) { */ if (!($interfaceIsEqual(err$1, $ifaceNil))) {} else { $s = 11; continue; }
							s.log(new sliceType$1([new $String("endpoint discovery:"), err$1]));
						/* } else { */ $s = 12; continue; case 11: 
							s.log(new sliceType$1([new $String("endpoint discovered")]));
							/* if ($pkg.WebSocketSupported && !s.forceLongPoll) { */ if ($pkg.WebSocketSupported && !s.forceLongPoll) {} else { $s = 13; continue; }
								_r$1 = s.connect(WebSocketTransport, hosts, backoff, $BLOCKING); /* */ $s = 14; case 14: if (_r$1 && _r$1.$blocking) { _r$1 = _r$1(); }
								/* if (_r$1) { */ if (_r$1) {} else { $s = 15; continue; }
									/* continue; */ $s = 1; continue;
								/* } */ case 15:
							/* } */ case 13:
							_r$2 = s.connect(LongPollTransport, hosts, backoff, $BLOCKING); /* */ $s = 16; case 16: if (_r$2 && _r$2.$blocking) { _r$2 = _r$2(); }
							_r$2;
						/* } */ case 12:
					/* } */ case 10:
				/* } else if (_selection[0] === 1) { */ $s = 7; continue; case 6: 
					return;
				/* } */ case 7:
			/* } */ case 4:
			delay = backoff.Failure(new Duration(0, 60000));
			/* if ((delay.$high > 0 || (delay.$high === 0 && delay.$low > 0))) { */ if ((delay.$high > 0 || (delay.$high === 0 && delay.$low > 0))) {} else { $s = 17; continue; }
				s.log(new sliceType$1([new $String("sleeping")]));
				s.connState("disconnected");
				$r = Sleep(delay, $BLOCKING); /* */ $s = 18; case 18: if ($r && $r.$blocking) { $r = $r(); }
			/* } */ case 17:
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } } catch(err) { $err = err; } finally { $deferFrames.pop(); if ($curGoroutine.asleep && !$jumpToDefer) { throw null; } $s = -1; $callDeferred($deferred, $err); } }; $blocking_discover.$blocking = true; return $blocking_discover;
	};
	Session.prototype.discover = function($b) { return this.$val.discover($b); };
	Session.Ptr.prototype.connect = function(transport, hosts, backoff, $b) {
		var $args = arguments, $r, $s = 0, $this = this, _i, _r, _ref, _tuple, connWorked, delay, gotOnline, host, s, transportWorked = false, trial;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_connect = function() { s: while (true) { switch ($s) { case 0:
		s = $this;
		trial = 0;
		/* while (trial < 2) { */ case 1: if(!(trial < 2)) { $s = 2; continue; }
			_ref = hosts;
			_i = 0;
			/* while (_i < _ref.$length) { */ case 3: if(!(_i < _ref.$length)) { $s = 4; continue; }
				host = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
				s.connState("connecting");
				_r = transport(s, host, $BLOCKING); /* */ $s = 5; case 5: if (_r && _r.$blocking) { _r = _r(); }
				_tuple = _r; connWorked = _tuple[0]; gotOnline = _tuple[1];
				if (connWorked) {
					transportWorked = true;
				}
				if (gotOnline) {
					backoff.Success();
					return transportWorked;
				}
				if (s.stopped) {
					return transportWorked;
				}
				delay = backoff.Failure(new Duration(0, 60000));
				/* if ((delay.$high > 0 || (delay.$high === 0 && delay.$low > 0))) { */ if ((delay.$high > 0 || (delay.$high === 0 && delay.$low > 0))) {} else { $s = 6; continue; }
					s.log(new sliceType$1([new $String("sleeping")]));
					s.connState("disconnected");
					$r = Sleep(delay, $BLOCKING); /* */ $s = 7; case 7: if ($r && $r.$blocking) { $r = $r(); }
				/* } */ case 6:
				_i++;
			/* } */ $s = 3; continue; case 4:
			trial = trial + (1) >> 0;
		/* } */ $s = 1; continue; case 2:
		return transportWorked;
		/* */ case -1: } return; } }; $blocking_connect.$blocking = true; return $blocking_connect;
	};
	Session.prototype.connect = function(transport, hosts, backoff, $b) { return this.$val.connect(transport, hosts, backoff, $b); };
	Session.Ptr.prototype.canLogin = function() {
		var _i, _i$1, _ref, _ref$1, key, key$1, s, value, value$1, value$2, value$3;
		s = this;
		value = s.sessionParams.access_key;
		if (!(value === undefined) && !(value === null)) {
			return true;
		}
		value$1 = s.sessionParams.user_id;
		if (!(value$1 === undefined) && !(value$1 === null)) {
			_ref = new sliceType(["user_auth", "master_sign"]);
			_i = 0;
			while (_i < _ref.$length) {
				key = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
				value$2 = s.sessionParams[$externalize(key, $String)];
				if (!(value$2 === undefined) && !(value$2 === null)) {
					return true;
				}
				_i++;
			}
			return false;
		}
		_ref$1 = new sliceType(["identity_type", "identity_name", "identity_auth"]);
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			key$1 = ((_i$1 < 0 || _i$1 >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i$1]);
			value$3 = s.sessionParams[$externalize(key$1, $String)];
			if (value$3 === undefined || value$3 === null) {
				return false;
			}
			_i$1++;
		}
		return true;
	};
	Session.prototype.canLogin = function() { return this.$val.canLogin(); };
	Session.Ptr.prototype.makeCreateSessionAction = function() {
		var header = null, s;
		s = this;
		header = s.sessionParams;
		header.action = $externalize("create_session", $String);
		return header;
	};
	Session.prototype.makeCreateSessionAction = function() { return this.$val.makeCreateSessionAction(); };
	Session.Ptr.prototype.makeResumeSessionAction = function(session) {
		var header = null, s;
		s = this;
		header = NewObject();
		header.action = $externalize("resume_session", $String);
		header.event_id = $externalize(s.receivedEventId, $Uint64);
		if (session) {
			header.session_id = s.sessionId;
		}
		s.sendEventAck = false;
		s.ackedEventId = s.receivedEventId;
		return header;
	};
	Session.prototype.makeResumeSessionAction = function(session) { return this.$val.makeResumeSessionAction(session); };
	Session.Ptr.prototype.handleSessionEvent = function(header) {
		var _i, _ref, _tuple, err, eventId, newValue, ok = false, param, s, sessionId, userAuth, userId;
		s = this;
		_tuple = GetSessionEventCredentials(header); userId = _tuple[0]; userAuth = _tuple[1]; sessionId = _tuple[2]; eventId = _tuple[3]; ok = _tuple[4]; err = _tuple[5];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			s.log(new sliceType$1([new $String("session creation:"), err]));
		}
		if (!jsInvoke("NinchatClient.Session onSessionEvent callback", s.onSessionEvent, new sliceType$1([header]))) {
			ok = false;
		}
		if (!ok) {
			s.sessionId = null;
			s.stopped = true;
			return ok;
		}
		s.sessionParams.user_id = userId;
		if (!(userAuth === null)) {
			s.sessionParams.user_auth = userAuth;
		}
		_ref = new sliceType(["identity_type", "identity_name", "identity_auth"]);
		_i = 0;
		while (_i < _ref.$length) {
			param = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			newValue = s.sessionParams[$externalize(param + "_new", $String)];
			if (!(newValue === undefined)) {
				s.sessionParams[$externalize(param, $String)] = newValue;
			}
			_i++;
		}
		delete s.sessionParams[$externalize("access_key", $String)];
		delete s.sessionParams[$externalize("master_sign", $String)];
		s.sessionId = sessionId;
		if (s.sendBuffer.$length === 0) {
			s.lastActionId = new $Uint64(0, 0);
		}
		s.sendEventAck = false;
		s.receivedEventId = eventId;
		s.ackedEventId = new $Uint64(0, 0);
		s.log(new sliceType$1([new $String("session created")]));
		ok = true;
		return ok;
	};
	Session.prototype.handleSessionEvent = function(header) { return this.$val.handleSessionEvent(header); };
	Session.Ptr.prototype.handleEvent = function(header, payload) {
		var _tuple, _tuple$1, _tuple$2, action, actionId = new $Uint64(0, 0), err, err$1, errorReason, errorType, eventId, i, lastReply, needsAck = false, ok = false, s, sessionLost = false, x, x$1, x$2, x$3, x$4;
		s = this;
		_tuple = GetEventAndActionId(header); eventId = _tuple[0]; actionId = _tuple[1]; err = _tuple[2];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			s.log(new sliceType$1([new $String("event:"), err]));
			return [actionId, sessionLost, needsAck, ok];
		}
		if ((eventId.$high > 0 || (eventId.$high === 0 && eventId.$low > 0))) {
			s.receivedEventId = eventId;
			if (!s.sendEventAck) {
				if ((x = (x$1 = s.receivedEventId, x$2 = s.ackedEventId, new $Uint64(x$1.$high - x$2.$high, x$1.$low - x$2.$low)), (x.$high > sessionEventAckWindow.$high || (x.$high === sessionEventAckWindow.$high && x.$low >= sessionEventAckWindow.$low)))) {
					s.sendAck();
				} else {
					needsAck = true;
				}
			}
		}
		if ((actionId.$high > 0 || (actionId.$high === 0 && actionId.$low > 0))) {
			i = sort.Search(s.numSent, (function(i) {
				var action, x$3, x$4;
				action = (x$3 = s.sendBuffer, ((i < 0 || i >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + i]));
				return (x$4 = action.Id, (x$4.$high > actionId.$high || (x$4.$high === actionId.$high && x$4.$low >= actionId.$low)));
			}));
			if (i < s.numSent) {
				action = (x$3 = s.sendBuffer, ((i < 0 || i >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + i]));
				if ((x$4 = action.Id, (x$4.$high === actionId.$high && x$4.$low === actionId.$low))) {
					_tuple$1 = IsEventLastReply(header, action); lastReply = _tuple$1[0]; err$1 = _tuple$1[1];
					if (!($interfaceIsEqual(err$1, $ifaceNil))) {
						s.log(new sliceType$1([new $String("event:"), err$1]));
						return [actionId, sessionLost, needsAck, ok];
					}
					if (!(action.Deferred === ptrType$2.nil)) {
						if (lastReply) {
							action.Deferred.Resolve(new sliceType$1([header, payload]));
						} else {
							action.Deferred.Notify(new sliceType$1([header, payload]));
						}
					}
					if (lastReply) {
						s.sendBuffer = $appendSlice($subslice(s.sendBuffer, 0, i), $subslice(s.sendBuffer, (i + 1 >> 0)));
						s.numSent = s.numSent - (1) >> 0;
					}
				}
			}
		}
		_tuple$2 = GetEventError(header); errorType = _tuple$2[0]; errorReason = _tuple$2[1]; sessionLost = _tuple$2[2]; err = _tuple$2[3];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			s.log(new sliceType$1([new $String("event:"), err]));
			if (sessionLost) {
				s.sessionId = null;
				if (!s.canLogin()) {
					jsInvoke("NinchatClient.Session onSessionEvent callback", s.onSessionEvent, new sliceType$1([header]));
					s.stopped = true;
				}
			}
			return [actionId, sessionLost, needsAck, ok];
		}
		if (errorType === "deprecated") {
			s.log(new sliceType$1([new $String("deprecated:"), new $String(errorReason)]));
		}
		if (!jsInvoke("NinchatClient.Session onEvent callback", s.onEvent, new sliceType$1([header, payload]))) {
			return [actionId, sessionLost, needsAck, ok];
		}
		ok = true;
		return [actionId, sessionLost, needsAck, ok];
	};
	Session.prototype.handleEvent = function(header, payload) { return this.$val.handleEvent(header, payload); };
	Session.Ptr.prototype.connState = function(state) {
		var s;
		s = this;
		if (!(s.latestConnState === state)) {
			s.latestConnState = state;
			if (!(s.onConnState === null)) {
				jsInvoke("NinchatClient.Session onConnState callback", s.onConnState, new sliceType$1([new $String(s.latestConnState)]));
			}
		}
	};
	Session.prototype.connState = function(state) { return this.$val.connState(state); };
	Session.Ptr.prototype.connActive = function() {
		var s;
		s = this;
		s.latestConnActive = Now();
		if (!(s.onConnActive === null)) {
			jsInvoke("NinchatClient.Session onConnActive callback", s.onConnActive, new sliceType$1([s.latestConnActive]));
		}
	};
	Session.prototype.connActive = function() { return this.$val.connActive(); };
	Session.Ptr.prototype.log = function(tokens) {
		var s;
		s = this;
		Log("NinchatClient.Session onLog callback", s.onLog, tokens);
	};
	Session.prototype.log = function(tokens) { return this.$val.log(tokens); };
	Now = $pkg.Now = function() {
		var x;
		return (x = $internalize(new ($global.Date)().getTime(), $Int64), new Time(x.$high, x.$low));
	};
	NewTimer = $pkg.NewTimer = function(timeout) {
		var timer = ptrType$4.nil;
		timer = new Timer.Ptr(new chanType$1(0), null);
		if ((timeout.$high > 0 || (timeout.$high === 0 && timeout.$low >= 0))) {
			timer.Reset(timeout);
		}
		return timer;
	};
	Timer.Ptr.prototype.Active = function() {
		var timer;
		timer = this;
		return !(timer.id === null);
	};
	Timer.prototype.Active = function() { return this.$val.Active(); };
	Timer.Ptr.prototype.Reset = function(timeout) {
		var timer;
		timer = this;
		timer.Stop();
		timer.id = SetTimeout((function() {
			timer.id = null;
			$go((function($b) {
				var $args = arguments, $r, $s = 0, $this = this;
				/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
				$r = $send(timer.C, true, $BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
				/* */ case -1: } return; } }; $f.$blocking = true; return $f;
			}), []);
		}), timeout);
	};
	Timer.prototype.Reset = function(timeout) { return this.$val.Reset(timeout); };
	Timer.Ptr.prototype.Stop = function() {
		var timer;
		timer = this;
		if (!(timer.id === null)) {
			ClearTimeout(timer.id);
			timer.id = null;
		}
	};
	Timer.prototype.Stop = function() { return this.$val.Stop(); };
	Sleep = $pkg.Sleep = function(delay, $b) {
		var $args = arguments, $r, $s = 0, $this = this, _r;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_Sleep = function() { s: while (true) { switch ($s) { case 0:
		_r = $recv(NewTimer(delay).C, $BLOCKING); /* */ $s = 1; case 1: if (_r && _r.$blocking) { _r = _r(); }
		_r[0];
		/* */ case -1: } return; } }; $blocking_Sleep.$blocking = true; return $blocking_Sleep;
	};
	NewWebSocket = $pkg.NewWebSocket = function(url) {
		var ws = ptrType$5.nil;
		ws = new WebSocket.Ptr(new chanType$1(1), new ($global.WebSocket)($externalize(url, $String)), false, $ifaceNil, sliceType$2.nil);
		ws.impl.binaryType = $externalize("arraybuffer", $String);
		ws.impl.onopen = $externalize((function(param) {
			ws.open = true;
			$go((function($b) {
				var $args = arguments, $r, $s = 0, $this = this;
				/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
				$r = $send(ws.Notify, true, $BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
				/* */ case -1: } return; } }; $f.$blocking = true; return $f;
			}), []);
		}), funcType$10);
		ws.impl.onmessage = $externalize((function(object) {
			ws.buffer = $append(ws.buffer, object.data);
			$go((function() {
				var _selection;
				_selection = $select([[ws.Notify, true], []]);
			}), []);
		}), funcType$11);
		ws.impl.onclose = $externalize((function(param) {
			ws.open = false;
			$go((function() {
				$close(ws.Notify);
			}), []);
		}), funcType$10);
		ws.impl.onerror = $externalize((function(object) {
			ws.error = errors.New("WebSocket error event");
		}), funcType$11);
		return ws;
	};
	WebSocket.Ptr.prototype.Send = function(data) {
		var $deferred = [], $err = null, err = $ifaceNil, ws;
		/* */ try { $deferFrames.push($deferred);
		ws = this;
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		err = ws.error;
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			return err;
		}
		ws.impl.send($externalize(data, $emptyInterface));
		return err;
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return err; }
	};
	WebSocket.prototype.Send = function(data) { return this.$val.Send(data); };
	WebSocket.Ptr.prototype.SendJSON = function(object) {
		var _tuple, err = $ifaceNil, json, ws;
		ws = this;
		_tuple = StringifyJSON(object); json = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			return err;
		}
		err = ws.Send(new $String(json));
		return err;
	};
	WebSocket.prototype.SendJSON = function(object) { return this.$val.SendJSON(object); };
	WebSocket.Ptr.prototype.Receive = function() {
		var data = null, err = $ifaceNil, ws, x;
		ws = this;
		err = ws.error;
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			return [data, err];
		}
		if (!ws.open) {
			return [data, err];
		}
		if (ws.buffer.$length > 0) {
			data = (x = ws.buffer, ((0 < 0 || 0 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + 0]));
			ws.buffer = $subslice(ws.buffer, 1);
		}
		return [data, err];
	};
	WebSocket.prototype.Receive = function() { return this.$val.Receive(); };
	WebSocket.Ptr.prototype.ReceiveJSON = function() {
		var _tuple, _tuple$1, data, err = $ifaceNil, object = null, ws;
		ws = this;
		_tuple = ws.Receive(); data = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, $ifaceNil)) || data === null) {
			return [object, err];
		}
		_tuple$1 = ParseJSON(StringifyFrame(data)); object = _tuple$1[0]; err = _tuple$1[1];
		return [object, err];
	};
	WebSocket.prototype.ReceiveJSON = function() { return this.$val.ReceiveJSON(); };
	WebSocket.Ptr.prototype.Close = function() {
		var $deferred = [], $err = null, err = $ifaceNil, ws;
		/* */ try { $deferFrames.push($deferred);
		ws = this;
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		ws.impl.close();
		err = ws.error;
		return err;
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return err; }
	};
	WebSocket.prototype.Close = function() { return this.$val.Close(); };
	StringifyFrame = $pkg.StringifyFrame = function(data) {
		var _tuple, bytes, ok, s = "", view;
		_tuple = $assertType($internalize(data, $emptyInterface), $String, true); s = _tuple[0]; ok = _tuple[1];
		if (ok) {
			return s;
		}
		view = NewUint8Array(data);
		bytes = $assertType($internalize(view, $emptyInterface), sliceType$4);
		s = $bytesToString(bytes);
		return s;
	};
	WebSocketTransport = $pkg.WebSocketTransport = function(s, host, $b) {
		var $args = arguments, $deferred = [], $err = null, $r, $s = 0, $this = this, _r, _r$1, _selection, _tuple, connWorked = false, connectTimer, connected, gotOnline = false, hostHealthy, ws;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_WebSocketTransport = function() { try { $deferFrames.push($deferred); s: while (true) { switch ($s) { case 0:
		ws = ptrType$5.nil;
		$deferred.push([(function() {
			if (!(ws === ptrType$5.nil)) {
				ws.Close();
			}
		}), [$BLOCKING]]);
		connectTimer = NewTimer(new Duration(-1, 4294967295));
		$deferred.push([$methodVal(connectTimer, "Stop"), [$BLOCKING]]);
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			gotOnline = false;
			hostHealthy = false;
			s.log(new sliceType$1([new $String("connecting to"), new $String(host)]));
			ws = NewWebSocket("wss://" + host + "/v2/socket");
			connectTimer.Reset(JitterDuration(new Duration(0, 9000), 0.1));
			_r = $select([[ws.Notify], [connectTimer.C], [s.closeNotify]], $BLOCKING); /* */ $s = 3; case 3: if (_r && _r.$blocking) { _r = _r(); }
			_selection = _r;
			/* if (_selection[0] === 0) { */ if (_selection[0] === 0) {} else if (_selection[0] === 1) { $s = 4; continue; } else if (_selection[0] === 2) { $s = 5; continue; } else { $s = 6; continue; }
				connected = _selection[1][0];
				connectTimer.Stop();
				/* if (connected) { */ if (connected) {} else { $s = 7; continue; }
					s.log(new sliceType$1([new $String("connected")]));
					s.connState("connected");
					connWorked = true;
					_r$1 = webSocketHandshake(s, ws, $BLOCKING); /* */ $s = 9; case 9: if (_r$1 && _r$1.$blocking) { _r$1 = _r$1(); }
					_tuple = _r$1; gotOnline = _tuple[0]; hostHealthy = _tuple[1];
				/* } else { */ $s = 8; continue; case 7: 
					s.log(new sliceType$1([new $String("connection failed")]));
				/* } */ case 8:
			/* } else if (_selection[0] === 1) { */ $s = 6; continue; case 4: 
				s.log(new sliceType$1([new $String("connection timeout")]));
			/* } else if (_selection[0] === 2) { */ $s = 6; continue; case 5: 
				connectTimer.Stop();
			/* } */ case 6:
			ws.Close();
			ws = ptrType$5.nil;
			s.log(new sliceType$1([new $String("disconnected")]));
			if (!gotOnline || !hostHealthy || s.stopped) {
				return [connWorked, gotOnline];
			}
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } } catch(err) { $err = err; } finally { $deferFrames.pop(); if ($curGoroutine.asleep && !$jumpToDefer) { throw null; } $s = -1; $callDeferred($deferred, $err); return [connWorked, gotOnline]; } }; $blocking_WebSocketTransport.$blocking = true; return $blocking_WebSocketTransport;
	};
	webSocketHandshake = function(s, ws, $b) {
		var $args = arguments, $r, $s = 0, $this = this, _r, _r$1, _r$2, _selection, _tuple, _tuple$1, connected, done, err, err$1, fail, gotEvents, gotOnline = false, header, header$1, hostHealthy = false, timer;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_webSocketHandshake = function() { s: while (true) { switch ($s) { case 0:
		header = null;
		if (s.sessionId === null) {
			s.log(new sliceType$1([new $String("session creation")]));
			header = s.makeCreateSessionAction();
		} else {
			s.log(new sliceType$1([new $String("session resumption")]));
			header = s.makeResumeSessionAction(true);
		}
		err = ws.SendJSON(header);
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			s.log(new sliceType$1([new $String("send:"), err]));
		}
		/* if (s.sessionId === null) { */ if (s.sessionId === null) {} else { $s = 1; continue; }
			header$1 = null;
			timer = NewTimer(JitterDuration(new Duration(0, 13000), 0.2));
			/* while (true) { */ case 2: if(!(true)) { $s = 3; continue; }
				err$1 = $ifaceNil;
				_tuple = ws.ReceiveJSON(); header$1 = _tuple[0]; err$1 = _tuple[1];
				if (!($interfaceIsEqual(err$1, $ifaceNil))) {
					s.log(new sliceType$1([new $String("session creation:"), err$1]));
					return [gotOnline, hostHealthy];
				}
				if (!(header$1 === null)) {
					timer.Stop();
					/* break; */ $s = 3; continue;
				}
				_r = $select([[ws.Notify], [timer.C]], $BLOCKING); /* */ $s = 4; case 4: if (_r && _r.$blocking) { _r = _r(); }
				_selection = _r;
				if (_selection[0] === 0) {
					connected = _selection[1][0];
					if (!connected) {
						s.log(new sliceType$1([new $String("disconnected during session creation")]));
						timer.Stop();
						return [gotOnline, hostHealthy];
					}
				} else if (_selection[0] === 1) {
					s.log(new sliceType$1([new $String("session creation timeout")]));
					return [gotOnline, hostHealthy];
				}
			/* } */ $s = 2; continue; case 3:
			if (!s.handleSessionEvent(header$1)) {
				return [gotOnline, hostHealthy];
			}
			gotOnline = true;
			hostHealthy = true;
			s.connActive();
		/* } */ case 1:
		fail = new chanType$1(1);
		done = new chanType$1(0);
		$go(webSocketSend, [s, ws, fail, done]);
		_r$1 = webSocketReceive(s, ws, fail, $BLOCKING); /* */ $s = 5; case 5: if (_r$1 && _r$1.$blocking) { _r$1 = _r$1(); }
		_tuple$1 = _r$1; gotEvents = _tuple$1[0]; hostHealthy = _tuple$1[1];
		if (gotEvents) {
			gotOnline = true;
		}
		_r$2 = $recv(done, $BLOCKING); /* */ $s = 6; case 6: if (_r$2 && _r$2.$blocking) { _r$2 = _r$2(); }
		_r$2[0];
		return [gotOnline, hostHealthy];
		/* */ case -1: } return; } }; $blocking_webSocketHandshake.$blocking = true; return $blocking_webSocketHandshake;
	};
	webSocketSend = function(s, ws, fail, done, $b) {
		var $args = arguments, $deferred = [], $err = null, $r, $s = 0, $this = this, _key, _map, _r, _selection, _tuple, _tuple$1, _tuple$2, action, action$1, array, base64, buffer, closeSession, data, err, err$1, err$2, err$3, err$4, err$5, frame, i, i$1, keeper, length, ok, sending, x, x$1, x$2, x$3, x$4, x$5, x$6;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_webSocketSend = function() { try { $deferFrames.push($deferred); s: while (true) { switch ($s) { case 0:
		$deferred.push([(function($b) {
			var $args = arguments, $r, $s = 0, $this = this;
			/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
			$r = $send(done, true, $BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
			/* */ case -1: } return; } }; $f.$blocking = true; return $f;
		}), [$BLOCKING]]);
		keeper = NewTimer(JitterDuration(new Duration(0, 56000), -0.3));
		$deferred.push([$methodVal(keeper, "Stop"), [$BLOCKING]]);
		s.numSent = 0;
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			/* while (s.numSent < s.sendBuffer.$length) { */ case 3: if(!(s.numSent < s.sendBuffer.$length)) { $s = 4; continue; }
				action = (x = s.sendBuffer, x$1 = s.numSent, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1]));
				if (!(action.Payload === null)) {
					action.Header.frames = $parseInt(action.Payload.length);
				}
				if (!((x$2 = s.receivedEventId, x$3 = s.ackedEventId, (x$2.$high === x$3.$high && x$2.$low === x$3.$low)))) {
					action.Header.event_id = $externalize(s.receivedEventId, $Uint64);
					s.sendEventAck = false;
					s.ackedEventId = s.receivedEventId;
				}
				err = ws.SendJSON(action.Header);
				delete action.Header[$externalize("frames", $String)];
				delete action.Header[$externalize("event_id", $String)];
				/* if (!($interfaceIsEqual(err, $ifaceNil))) { */ if (!($interfaceIsEqual(err, $ifaceNil))) {} else { $s = 5; continue; }
					s.log(new sliceType$1([new $String("send:"), err]));
					$r = $send(fail, true, $BLOCKING); /* */ $s = 6; case 6: if ($r && $r.$blocking) { $r = $r(); }
					return;
				/* } */ case 5:
				/* if (!(action.Payload === null)) { */ if (!(action.Payload === null)) {} else { $s = 7; continue; }
					i = 0;
					/* while (i < $parseInt(action.Payload.length)) { */ case 8: if(!(i < $parseInt(action.Payload.length))) { $s = 9; continue; }
						frame = action.Payload[i];
						/* if ($internalize(action.Header.action, $String) === "update_user") { */ if ($internalize(action.Header.action, $String) === "update_user") {} else { $s = 10; continue; }
							_tuple = $assertType($internalize(frame, $emptyInterface), $String, true); ok = _tuple[1];
							/* if (ok) { */ if (ok) {} else { $s = 11; continue; }
								_tuple$1 = ParseDataURI(frame); base64 = _tuple$1[0]; err$1 = _tuple$1[1];
								/* if (!($interfaceIsEqual(err$1, $ifaceNil))) { */ if (!($interfaceIsEqual(err$1, $ifaceNil))) {} else { $s = 12; continue; }
									s.log(new sliceType$1([new $String("send:"), err$1]));
									$r = $send(fail, true, $BLOCKING); /* */ $s = 13; case 13: if ($r && $r.$blocking) { $r = $r(); }
									return;
								/* } */ case 12:
								_tuple$2 = Atob(base64); data = _tuple$2[0]; err$1 = _tuple$2[1];
								/* if (!($interfaceIsEqual(err$1, $ifaceNil))) { */ if (!($interfaceIsEqual(err$1, $ifaceNil))) {} else { $s = 14; continue; }
									s.log(new sliceType$1([new $String("send:"), err$1]));
									$r = $send(fail, true, $BLOCKING); /* */ $s = 15; case 15: if ($r && $r.$blocking) { $r = $r(); }
									return;
								/* } */ case 14:
								length = $parseInt(data.length);
								buffer = NewArrayBuffer(length);
								array = NewUint8Array(buffer);
								i$1 = 0;
								while (i$1 < length) {
									array[i$1] = data.charCodeAt(i$1);
									i$1 = i$1 + (1) >> 0;
								}
								frame = buffer;
							/* } */ case 11:
						/* } */ case 10:
						err$2 = ws.Send(frame);
						/* if (!($interfaceIsEqual(err$2, $ifaceNil))) { */ if (!($interfaceIsEqual(err$2, $ifaceNil))) {} else { $s = 16; continue; }
							s.log(new sliceType$1([new $String("send:"), err$2]));
							$r = $send(fail, true, $BLOCKING); /* */ $s = 17; case 17: if ($r && $r.$blocking) { $r = $r(); }
							return;
						/* } */ case 16:
						i = i + (1) >> 0;
					/* } */ $s = 8; continue; case 9:
				/* } */ case 7:
				if ((x$4 = action.Id, (x$4.$high === 0 && x$4.$low === 0))) {
					s.sendBuffer = $appendSlice($subslice(s.sendBuffer, 0, s.numSent), $subslice(s.sendBuffer, (s.numSent + 1 >> 0)));
				} else {
					s.numSent = s.numSent + (1) >> 0;
				}
				keeper.Reset(JitterDuration(new Duration(0, 56000), -0.3));
			/* } */ $s = 3; continue; case 4:
			/* if (s.sendEventAck && !((x$5 = s.receivedEventId, x$6 = s.ackedEventId, (x$5.$high === x$6.$high && x$5.$low === x$6.$low)))) { */ if (s.sendEventAck && !((x$5 = s.receivedEventId, x$6 = s.ackedEventId, (x$5.$high === x$6.$high && x$5.$low === x$6.$low)))) {} else { $s = 18; continue; }
				action$1 = s.makeResumeSessionAction(false);
				err$3 = ws.SendJSON(action$1);
				/* if (!($interfaceIsEqual(err$3, $ifaceNil))) { */ if (!($interfaceIsEqual(err$3, $ifaceNil))) {} else { $s = 19; continue; }
					s.log(new sliceType$1([new $String("send:"), err$3]));
					$r = $send(fail, true, $BLOCKING); /* */ $s = 20; case 20: if ($r && $r.$blocking) { $r = $r(); }
					return;
				/* } */ case 19:
			/* } */ case 18:
			_r = $select([[s.sendNotify], [keeper.C], [fail]], $BLOCKING); /* */ $s = 21; case 21: if (_r && _r.$blocking) { _r = _r(); }
			_selection = _r;
			/* if (_selection[0] === 0) { */ if (_selection[0] === 0) {} else if (_selection[0] === 1) { $s = 22; continue; } else if (_selection[0] === 2) { $s = 23; continue; } else { $s = 24; continue; }
				sending = _selection[1][0];
				if (!sending) {
					closeSession = (_map = new $Map(), _key = "action", _map[_key] = { k: _key, v: new $String("close_session") }, _map);
					err$4 = ws.SendJSON(new mapType(closeSession));
					if (!($interfaceIsEqual(err$4, $ifaceNil))) {
						s.log(new sliceType$1([new $String("send:"), err$4]));
					}
					return;
				}
			/* } else if (_selection[0] === 1) { */ $s = 24; continue; case 22: 
				err$5 = ws.Send(new sliceType$5([]));
				/* if (!($interfaceIsEqual(err$5, $ifaceNil))) { */ if (!($interfaceIsEqual(err$5, $ifaceNil))) {} else { $s = 25; continue; }
					s.log(new sliceType$1([new $String("send:"), err$5]));
					$r = $send(fail, true, $BLOCKING); /* */ $s = 26; case 26: if ($r && $r.$blocking) { $r = $r(); }
					return;
				/* } */ case 25:
				keeper.Reset(JitterDuration(new Duration(0, 56000), -0.3));
			/* } else if (_selection[0] === 2) { */ $s = 24; continue; case 23: 
				return;
			/* } */ case 24:
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } } catch(err) { $err = err; } finally { $deferFrames.pop(); if ($curGoroutine.asleep && !$jumpToDefer) { throw null; } $s = -1; $callDeferred($deferred, $err); } }; $blocking_webSocketSend.$blocking = true; return $blocking_webSocketSend;
	};
	webSocketReceive = function(s, ws, fail, $b) {
		var $args = arguments, $deferred = [], $err = null, $r, $s = 0, $this = this, _r, _selection, _selection$1, _tuple, _tuple$1, _tuple$2, _tuple$3, _tuple$4, ackNeeded, acker, connected, data, data$1, err, err$1, frames, gotEvents = false, header, hostHealthy = false, needsAck, ok, payload, sessionLost, text, watchdog, x, x$1;
		/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $blocking_webSocketReceive = function() { try { $deferFrames.push($deferred); s: while (true) { switch ($s) { case 0:
		header = null;
		payload = null;
		frames = 0;
		watchdog = NewTimer(JitterDuration(new Duration(0, 64000), 0.3));
		$deferred.push([$methodVal(watchdog, "Stop"), [$BLOCKING]]);
		acker = NewTimer(new Duration(-1, 4294967295));
		$deferred.push([$methodVal(acker, "Stop"), [$BLOCKING]]);
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			ackNeeded = false;
			/* while (true) { */ case 3: if(!(true)) { $s = 4; continue; }
				/* if (header === null) { */ if (header === null) {} else { $s = 5; continue; }
					err = $ifaceNil;
					_tuple = ws.Receive(); data = _tuple[0]; err = _tuple[1];
					/* if (!($interfaceIsEqual(err, $ifaceNil))) { */ if (!($interfaceIsEqual(err, $ifaceNil))) {} else { $s = 7; continue; }
						s.log(new sliceType$1([new $String("receive:"), err]));
						hostHealthy = false;
						$r = $send(fail, true, $BLOCKING); /* */ $s = 8; case 8: if ($r && $r.$blocking) { $r = $r(); }
						return [gotEvents, hostHealthy];
					/* } */ case 7:
					if (data === null) {
						/* break; */ $s = 4; continue;
					}
					watchdog.Reset(JitterDuration(new Duration(0, 64000), 0.7));
					s.connActive();
					text = StringifyFrame(data);
					if (text.length === 0) {
						/* continue; */ $s = 3; continue;
					}
					_tuple$1 = ParseJSON(text); header = _tuple$1[0]; err = _tuple$1[1];
					/* if (!($interfaceIsEqual(err, $ifaceNil))) { */ if (!($interfaceIsEqual(err, $ifaceNil))) {} else { $s = 9; continue; }
						s.log(new sliceType$1([new $String("receive:"), err]));
						hostHealthy = false;
						$r = $send(fail, true, $BLOCKING); /* */ $s = 10; case 10: if ($r && $r.$blocking) { $r = $r(); }
						return [gotEvents, hostHealthy];
					/* } */ case 9:
					payload = NewArray();
					_tuple$2 = GetEventFrames(header); frames = _tuple$2[0]; err = _tuple$2[1];
					/* if (!($interfaceIsEqual(err, $ifaceNil))) { */ if (!($interfaceIsEqual(err, $ifaceNil))) {} else { $s = 11; continue; }
						s.log(new sliceType$1([new $String("receive:"), err]));
						hostHealthy = false;
						$r = $send(fail, true, $BLOCKING); /* */ $s = 12; case 12: if ($r && $r.$blocking) { $r = $r(); }
						return [gotEvents, hostHealthy];
					/* } */ case 11:
				/* } else { */ $s = 6; continue; case 5: 
					_tuple$3 = ws.Receive(); data$1 = _tuple$3[0]; err$1 = _tuple$3[1];
					/* if (!($interfaceIsEqual(err$1, $ifaceNil))) { */ if (!($interfaceIsEqual(err$1, $ifaceNil))) {} else { $s = 13; continue; }
						s.log(new sliceType$1([new $String("receive:"), err$1]));
						hostHealthy = false;
						$r = $send(fail, true, $BLOCKING); /* */ $s = 14; case 14: if ($r && $r.$blocking) { $r = $r(); }
						return [gotEvents, hostHealthy];
					/* } */ case 13:
					if (data$1 === null) {
						/* break; */ $s = 4; continue;
					}
					payload.push(data$1);
					frames = frames - (1) >> 0;
				/* } */ case 6:
				/* if (frames === 0) { */ if (frames === 0) {} else { $s = 15; continue; }
					_tuple$4 = s.handleEvent(header, payload); sessionLost = _tuple$4[1]; needsAck = _tuple$4[2]; ok = _tuple$4[3];
					/* if (!ok) { */ if (!ok) {} else { $s = 16; continue; }
						if (sessionLost) {
							gotEvents = true;
						} else {
							hostHealthy = false;
						}
						$r = $send(fail, true, $BLOCKING); /* */ $s = 17; case 17: if ($r && $r.$blocking) { $r = $r(); }
						return [gotEvents, hostHealthy];
					/* } */ case 16:
					if (needsAck) {
						ackNeeded = true;
					}
					header = null;
					payload = null;
					frames = 0;
					gotEvents = true;
					hostHealthy = true;
				/* } */ case 15:
				_selection = $select([[s.closeNotify], [fail], []]);
				if (_selection[0] === 0) {
					return [gotEvents, hostHealthy];
				} else if (_selection[0] === 1) {
					return [gotEvents, hostHealthy];
				}
			/* } */ $s = 3; continue; case 4:
			if (ackNeeded && !acker.Active()) {
				acker.Reset(JitterDuration(new Duration(0, 7000), -0.3));
			}
			_r = $select([[ws.Notify], [watchdog.C], [acker.C], [s.closeNotify], [fail]], $BLOCKING); /* */ $s = 18; case 18: if (_r && _r.$blocking) { _r = _r(); }
			_selection$1 = _r;
			/* if (_selection$1[0] === 0) { */ if (_selection$1[0] === 0) {} else if (_selection$1[0] === 1) { $s = 19; continue; } else if (_selection$1[0] === 2) { $s = 20; continue; } else if (_selection$1[0] === 3) { $s = 21; continue; } else if (_selection$1[0] === 4) { $s = 22; continue; } else { $s = 23; continue; }
				connected = _selection$1[1][0];
				/* if (!connected) { */ if (!connected) {} else { $s = 24; continue; }
					$r = $send(fail, true, $BLOCKING); /* */ $s = 25; case 25: if ($r && $r.$blocking) { $r = $r(); }
					return [gotEvents, hostHealthy];
				/* } */ case 24:
			/* } else if (_selection$1[0] === 1) { */ $s = 23; continue; case 19: 
				s.log(new sliceType$1([new $String("receive timeout")]));
				$r = $send(fail, true, $BLOCKING); /* */ $s = 26; case 26: if ($r && $r.$blocking) { $r = $r(); }
				return [gotEvents, hostHealthy];
			/* } else if (_selection$1[0] === 2) { */ $s = 23; continue; case 20: 
				if (!s.sendEventAck && !((x = s.ackedEventId, x$1 = s.receivedEventId, (x.$high === x$1.$high && x.$low === x$1.$low)))) {
					s.sendAck();
				}
			/* } else if (_selection$1[0] === 3) { */ $s = 23; continue; case 21: 
				return [gotEvents, hostHealthy];
			/* } else if (_selection$1[0] === 4) { */ $s = 23; continue; case 22: 
				return [gotEvents, hostHealthy];
			/* } */ case 23:
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } } catch(err) { $err = err; } finally { $deferFrames.pop(); if ($curGoroutine.asleep && !$jumpToDefer) { throw null; } $s = -1; $callDeferred($deferred, $err); return [gotEvents, hostHealthy]; } }; $blocking_webSocketReceive.$blocking = true; return $blocking_webSocketReceive;
	};
	init = function() {
		xhrType = $global.XDomainRequest;
		if (xhrType === undefined) {
			xhrType = $global.XMLHttpRequest;
			xhrRequestHeaderSupport = true;
		}
	};
	XHR = $pkg.XHR = function(url, data, timeout) {
		var $deferred = [], $err = null, channel = chanType$2.nil, err = $ifaceNil, method, request;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			err = jsError($recover());
		}), []]);
		method = "";
		if (data === "") {
			method = "GET";
		} else {
			method = "POST";
		}
		channel = new chanType$2(1);
		request = new (xhrType)();
		request.onload = $externalize((function() {
			var obj, response;
			response = "";
			obj = request.responseText;
			if (!(obj === undefined) && !(obj === null)) {
				response = $internalize(obj, $String);
			}
			$go((function($b) {
				var $args = arguments, $r, $s = 0, $this = this;
				/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
				$r = $send(channel, response, $BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
				/* */ case -1: } return; } }; $f.$blocking = true; return $f;
			}), []);
		}), funcType);
		request.onprogress = $externalize((function() {
		}), funcType);
		request.ontimeout = $externalize((function() {
			$go((function() {
				$close(channel);
			}), []);
		}), funcType);
		request.onerror = $externalize((function() {
			$go((function($b) {
				var $args = arguments, $r, $s = 0, $this = this;
				/* */ if($b !== $BLOCKING) { $nonblockingCall(); }; var $f = function() { s: while (true) { switch ($s) { case 0:
				$r = $send(channel, "", $BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
				/* */ case -1: } return; } }; $f.$blocking = true; return $f;
			}), []);
		}), funcType);
		request.open($externalize(method, $String), $externalize(url, $String));
		request.timeout = $externalize(timeout, Duration);
		if (!(data === "") && xhrRequestHeaderSupport) {
			request.setRequestHeader($externalize("Content-Type", $String), $externalize("application/json", $String));
		}
		request.send($externalize(data, $String));
		return [channel, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [channel, err]; }
	};
	XHR_JSON = $pkg.XHR_JSON = function(url, data, timeout) {
		var _tuple, _tuple$1, channel = chanType$2.nil, err = $ifaceNil, json;
		_tuple = StringifyJSON(data); json = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, $ifaceNil))) {
			return [channel, err];
		}
		_tuple$1 = XHR(url, json, timeout); channel = _tuple$1[0]; err = _tuple$1[1];
		return [channel, err];
	};
	$pkg.$init = function() {
		$pkg.$init = function() {};
		/* */ var $r, $s = 0; var $init_main = function() { while (true) { switch ($s) { case 0:
		$r = errors.$init($BLOCKING); /* */ $s = 1; case 1: if ($r && $r.$blocking) { $r = $r(); }
		$r = js.$init($BLOCKING); /* */ $s = 2; case 2: if ($r && $r.$blocking) { $r = $r(); }
		$r = sort.$init($BLOCKING); /* */ $s = 3; case 3: if ($r && $r.$blocking) { $r = $r(); }
		sliceType = $sliceType($String);
		mapType = $mapType($String, $emptyInterface);
		sliceType$1 = $sliceType($emptyInterface);
		ptrType = $ptrType(js.Error);
		funcType = $funcType([], [], false);
		ptrType$1 = $ptrType($Bool);
		chanType = $chanType($String, false, true);
		funcType$1 = $funcType([js.Object, js.Object, js.Object], [$emptyInterface], false);
		funcType$2 = $funcType([], [mapType], false);
		funcType$3 = $funcType([js.Object], [$String], false);
		ptrType$2 = $ptrType(Deferred);
		sliceType$2 = $sliceType(js.Object);
		funcType$4 = $funcType([js.Object, js.Object, js.Object], [], false);
		chanType$1 = $chanType($Bool, false, false);
		ptrType$3 = $ptrType(Action);
		sliceType$3 = $sliceType(ptrType$3);
		funcType$5 = $funcType([js.Object], [], false);
		funcType$6 = $funcType([js.Object], [], false);
		funcType$7 = $funcType([js.Object], [], false);
		funcType$8 = $funcType([js.Object], [], false);
		funcType$9 = $funcType([js.Object, js.Object], [$emptyInterface], false);
		ptrType$4 = $ptrType(Timer);
		ptrType$5 = $ptrType(WebSocket);
		funcType$10 = $funcType([js.Object], [], false);
		funcType$11 = $funcType([js.Object], [], false);
		sliceType$4 = $sliceType($Uint8);
		sliceType$5 = $sliceType($Uint8);
		chanType$2 = $chanType($String, false, false);
		ptrType$6 = $ptrType(Backoff);
		ptrType$7 = $ptrType(Session);
		ptrType$3.methods = [["Name", "Name", "", $funcType([], [$String], false), -1]];
		Action.init([["Id", "Id", "", $Uint64, ""], ["Header", "Header", "", js.Object, ""], ["Payload", "Payload", "", js.Object, ""], ["Deferred", "Deferred", "", ptrType$2, ""], ["name", "name", "ninchatclient", $String, ""]]);
		ptrType$6.methods = [["Failure", "Failure", "", $funcType([Duration], [Duration], false), -1], ["Success", "Success", "", $funcType([], [], false), -1]];
		Backoff.init([["lastSlot", "lastSlot", "ninchatclient", $Int, ""]]);
		ptrType$2.methods = [["Notify", "Notify", "", $funcType([sliceType$1], [], true), -1], ["Reject", "Reject", "", $funcType([sliceType$1], [], true), -1], ["Resolve", "Resolve", "", $funcType([sliceType$1], [], true), -1], ["then", "then", "ninchatclient", $funcType([js.Object, js.Object, js.Object], [], false), -1]];
		Deferred.init([["resolve", "resolve", "ninchatclient", sliceType$2, ""], ["reject", "reject", "ninchatclient", sliceType$2, ""], ["notify", "notify", "ninchatclient", sliceType$2, ""]]);
		Transport.init([ptrType$7, $String], [$Bool, $Bool], false);
		ptrType$7.methods = [["Close", "Close", "", $funcType([], [], false), -1], ["OnConnActive", "OnConnActive", "", $funcType([js.Object], [], false), -1], ["OnConnState", "OnConnState", "", $funcType([js.Object], [], false), -1], ["OnEvent", "OnEvent", "", $funcType([js.Object], [], false), -1], ["OnLog", "OnLog", "", $funcType([js.Object], [], false), -1], ["OnSessionEvent", "OnSessionEvent", "", $funcType([js.Object], [], false), -1], ["Open", "Open", "", $funcType([], [], false), -1], ["Send", "Send", "", $funcType([js.Object, js.Object], [$emptyInterface], false), -1], ["SetAddress", "SetAddress", "", $funcType([js.Object], [], false), -1], ["SetParams", "SetParams", "", $funcType([js.Object], [], false), -1], ["SetTransport", "SetTransport", "", $funcType([js.Object], [], false), -1], ["canLogin", "canLogin", "ninchatclient", $funcType([], [$Bool], false), -1], ["connActive", "connActive", "ninchatclient", $funcType([], [], false), -1], ["connState", "connState", "ninchatclient", $funcType([$String], [], false), -1], ["connect", "connect", "ninchatclient", $funcType([Transport, sliceType, ptrType$6], [$Bool], false), -1], ["discover", "discover", "ninchatclient", $funcType([], [], false), -1], ["handleEvent", "handleEvent", "ninchatclient", $funcType([js.Object, js.Object], [$Uint64, $Bool, $Bool, $Bool], false), -1], ["handleSessionEvent", "handleSessionEvent", "ninchatclient", $funcType([js.Object], [$Bool], false), -1], ["log", "log", "ninchatclient", $funcType([sliceType$1], [], true), -1], ["makeCreateSessionAction", "makeCreateSessionAction", "ninchatclient", $funcType([], [js.Object], false), -1], ["makeResumeSessionAction", "makeResumeSessionAction", "ninchatclient", $funcType([$Bool], [js.Object], false), -1], ["send", "send", "ninchatclient", $funcType([ptrType$3], [], false), -1], ["sendAck", "sendAck", "ninchatclient", $funcType([], [], false), -1]];
		Session.init([["onSessionEvent", "onSessionEvent", "ninchatclient", js.Object, ""], ["onEvent", "onEvent", "ninchatclient", js.Object, ""], ["onConnState", "onConnState", "ninchatclient", js.Object, ""], ["onConnActive", "onConnActive", "ninchatclient", js.Object, ""], ["onLog", "onLog", "ninchatclient", js.Object, ""], ["address", "address", "ninchatclient", $String, ""], ["forceLongPoll", "forceLongPoll", "ninchatclient", $Bool, ""], ["sessionParams", "sessionParams", "ninchatclient", js.Object, ""], ["sessionId", "sessionId", "ninchatclient", js.Object, ""], ["latestConnState", "latestConnState", "ninchatclient", $String, ""], ["latestConnActive", "latestConnActive", "ninchatclient", Time, ""], ["lastActionId", "lastActionId", "ninchatclient", $Uint64, ""], ["sendNotify", "sendNotify", "ninchatclient", chanType$1, ""], ["sendBuffer", "sendBuffer", "ninchatclient", sliceType$3, ""], ["numSent", "numSent", "ninchatclient", $Int, ""], ["sendEventAck", "sendEventAck", "ninchatclient", $Bool, ""], ["receivedEventId", "receivedEventId", "ninchatclient", $Uint64, ""], ["ackedEventId", "ackedEventId", "ninchatclient", $Uint64, ""], ["closeNotify", "closeNotify", "ninchatclient", chanType$1, ""], ["closed", "closed", "ninchatclient", $Bool, ""], ["stopped", "stopped", "ninchatclient", $Bool, ""]]);
		ptrType$4.methods = [["Active", "Active", "", $funcType([], [$Bool], false), -1], ["Reset", "Reset", "", $funcType([Duration], [], false), -1], ["Stop", "Stop", "", $funcType([], [], false), -1]];
		Timer.init([["C", "C", "", chanType$1, ""], ["id", "id", "ninchatclient", js.Object, ""]]);
		ptrType$5.methods = [["Close", "Close", "", $funcType([], [$error], false), -1], ["Receive", "Receive", "", $funcType([], [js.Object, $error], false), -1], ["ReceiveJSON", "ReceiveJSON", "", $funcType([], [js.Object, $error], false), -1], ["Send", "Send", "", $funcType([$emptyInterface], [$error], false), -1], ["SendJSON", "SendJSON", "", $funcType([$emptyInterface], [$error], false), -1]];
		WebSocket.init([["Notify", "Notify", "", chanType$1, ""], ["impl", "impl", "ninchatclient", js.Object, ""], ["open", "open", "ninchatclient", $Bool, ""], ["error", "error", "ninchatclient", $error, ""], ["buffer", "buffer", "ninchatclient", sliceType$2, ""]]);
		xhrType = null;
		xhrRequestHeaderSupport = false;
		module = NewObject();
		sessionEventAckWindow = JitterUint64(new $Uint64(0, 4096), -0.25);
		$pkg.WebSocketSupported = !($global.WebSocket === undefined);
		init();
		main();
		/* */ } return; } }; $init_main.$blocking = true; return $init_main;
	};
	return $pkg;
})();
$packages["runtime"].$init()();
$go($packages["ninchatclient"].$init, [], true);
$flushConsole();

})(this);
//# sourceMappingURL=ninchatclient.js.map

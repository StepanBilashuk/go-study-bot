📖 Encoding & schema evolution (system design)

Data outlives code. You need serialization formats that let old and new versions coexist during rolling deploys.

▸ Formats
JSON/XML (verbose, schema-optional) vs binary Protobuf/Thrift/Avro (compact, schema-driven, field tags). Avro pairs data with a writer schema (great for data lakes / Kafka + a schema registry).

▸ Compatibility (define both with an example)
- Backward: new code reads OLD data (add fields as optional with defaults).
- Forward: old code reads NEW data (ignore unknown fields).
Rule: add optional fields, never reuse or renumber tags, don't remove required fields.

▸ Pitfalls
Changing a field's type, reusing a tag number, making a new field required.

▸ Interview probes
Forward vs backward compatibility with a concrete example; why Avro for pipelines; role of a schema registry.

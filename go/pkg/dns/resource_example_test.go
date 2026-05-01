// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleRawMessage_MarshalJSON() {
	core.Println(coverageExample(new(RawMessage).MarshalJSON()))
	// Output: 2
}

func ExampleRawMessage_UnmarshalJSON() {
	core.Println(coverageExample(new(RawMessage).UnmarshalJSON([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleResource_GetSize() {
	core.Println(coverageExample(new(Resource).GetSize()))
	// Output: 1
}

func ExampleNewResource() {
	core.Println(coverageExample(NewResource()))
	// Output: 1
}

func ExampleGetNewResource() {
	core.Println(coverageExample(GetNewResource()))
	// Output: 1
}

func ExampleDSRecord_Type() {
	core.Println(coverageExample(new(DSRecord).Type()))
	// Output: 1
}

func ExampleDSRecord_GetType() {
	core.Println(coverageExample(new(DSRecord).GetType()))
	// Output: 1
}

func ExampleNSRecord_Type() {
	core.Println(coverageExample(new(NSRecord).Type()))
	// Output: 1
}

func ExampleNSRecord_GetType() {
	core.Println(coverageExample(new(NSRecord).GetType()))
	// Output: 1
}

func ExampleNSRecord_GetNS() {
	core.Println(coverageExample(new(NSRecord).GetNS()))
	// Output: 1
}

func ExampleGLUE4Record_Type() {
	core.Println(coverageExample(new(GLUE4Record).Type()))
	// Output: 1
}

func ExampleGLUE4Record_GetType() {
	core.Println(coverageExample(new(GLUE4Record).GetType()))
	// Output: 1
}

func ExampleGLUE4Record_GetNS() {
	core.Println(coverageExample(new(GLUE4Record).GetNS()))
	// Output: 1
}

func ExampleGLUE6Record_Type() {
	core.Println(coverageExample(new(GLUE6Record).Type()))
	// Output: 1
}

func ExampleGLUE6Record_GetType() {
	core.Println(coverageExample(new(GLUE6Record).GetType()))
	// Output: 1
}

func ExampleGLUE6Record_GetNS() {
	core.Println(coverageExample(new(GLUE6Record).GetNS()))
	// Output: 1
}

func ExampleSYNTH4Record_Type() {
	core.Println(coverageExample(new(SYNTH4Record).Type()))
	// Output: 1
}

func ExampleSYNTH4Record_GetType() {
	core.Println(coverageExample(new(SYNTH4Record).GetType()))
	// Output: 1
}

func ExampleSYNTH4Record_GetNS() {
	core.Println(coverageExample(new(SYNTH4Record).GetNS()))
	// Output: 1
}

func ExampleSYNTH6Record_Type() {
	core.Println(coverageExample(new(SYNTH6Record).Type()))
	// Output: 1
}

func ExampleSYNTH6Record_GetType() {
	core.Println(coverageExample(new(SYNTH6Record).GetType()))
	// Output: 1
}

func ExampleSYNTH6Record_GetNS() {
	core.Println(coverageExample(new(SYNTH6Record).GetNS()))
	// Output: 1
}

func ExampleTXTRecord_Type() {
	core.Println(coverageExample(new(TXTRecord).Type()))
	// Output: 1
}

func ExampleTXTRecord_GetType() {
	core.Println(coverageExample(new(TXTRecord).GetType()))
	// Output: 1
}

func ExampleSYNTH4Record_NS() {
	core.Println(coverageExample(new(SYNTH4Record).NS()))
	// Output: 1
}

func ExampleSYNTH6Record_NS() {
	core.Println(coverageExample(new(SYNTH6Record).NS()))
	// Output: 1
}

func ExampleResource_HasType() {
	core.Println(coverageExample(new(Resource).HasType(1)))
	// Output: 1
}

func ExampleResource_GetHasType() {
	core.Println(coverageExample(new(Resource).GetHasType(1)))
	// Output: 1
}

func ExampleResource_HasNS() {
	core.Println(coverageExample(new(Resource).HasNS()))
	// Output: 1
}

func ExampleResource_GetHasNS() {
	core.Println(coverageExample(new(Resource).GetHasNS()))
	// Output: 1
}

func ExampleResource_HasDS() {
	core.Println(coverageExample(new(Resource).HasDS()))
	// Output: 1
}

func ExampleResource_GetHasDS() {
	core.Println(coverageExample(new(Resource).GetHasDS()))
	// Output: 1
}

func ExampleResource_ToNS() {
	core.Println(coverageExample(new(Resource).ToNS("lethean.lthn")))
	// Output: 1
}

func ExampleResource_GetToNS() {
	core.Println(coverageExample(new(Resource).GetToNS("lethean.lthn")))
	// Output: 1
}

func ExampleResource_ToGlue() {
	core.Println(coverageExample(new(Resource).ToGlue("lethean.lthn")))
	// Output: 1
}

func ExampleResource_GetToGlue() {
	core.Println(coverageExample(new(Resource).GetToGlue("lethean.lthn")))
	// Output: 1
}

func ExampleResource_ToDS() {
	core.Println(coverageExample(new(Resource).ToDS("lethean.lthn")))
	// Output: 1
}

func ExampleResource_GetToDS() {
	core.Println(coverageExample(new(Resource).GetToDS("lethean.lthn")))
	// Output: 1
}

func ExampleResource_ToTXT() {
	core.Println(coverageExample(new(Resource).ToTXT("lethean.lthn")))
	// Output: 1
}

func ExampleResource_GetToTXT() {
	core.Println(coverageExample(new(Resource).GetToTXT("lethean.lthn")))
	// Output: 1
}

func ExampleResource_ToZone() {
	core.Println(coverageExample(new(Resource).ToZone("lethean.lthn")))
	// Output: 1
}

func ExampleResource_GetToZone() {
	core.Println(coverageExample(new(Resource).GetToZone("lethean.lthn")))
	// Output: 1
}

func ExampleResource_ToReferral() {
	core.Println(coverageExample(new(Resource).ToReferral("lethean.lthn", 1, true)))
	// Output: 1
}

func ExampleResource_GetToReferral() {
	core.Println(coverageExample(new(Resource).GetToReferral("lethean.lthn", 1, true)))
	// Output: 1
}

func ExampleResource_ToDNS() {
	core.Println(coverageExample(new(Resource).ToDNS("lethean.lthn", 1)))
	// Output: 1
}

func ExampleResource_GetToDNS() {
	core.Println(coverageExample(new(Resource).GetToDNS("lethean.lthn", 1)))
	// Output: 1
}

func ExampleResource_ToNSEC() {
	core.Println(coverageExample(new(Resource).ToNSEC("lethean.lthn")))
	// Output: 1
}

func ExampleResource_GetToNSEC() {
	core.Println(coverageExample(new(Resource).GetToNSEC("lethean.lthn")))
	// Output: 1
}

func ExampleResource_Encode() {
	core.Println(coverageExample(new(Resource).Encode()))
	// Output: 2
}

func ExampleResource_GetEncode() {
	core.Println(coverageExample(new(Resource).GetEncode()))
	// Output: 2
}

func ExampleResource_Decode() {
	core.Println(coverageExample(new(Resource).Decode([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleResource_GetDecode() {
	core.Println(coverageExample(new(Resource).GetDecode([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleResource_GetJSON() {
	core.Println(coverageExample(new(Resource).GetJSON()))
	// Output: 1
}

func ExampleResource_FromJSON() {
	core.Println(coverageExample(new(Resource).FromJSON(coverageValue[ResourceJSON]("good"))))
	// Output: 1
}

func ExampleDSRecord_GetSize() {
	core.Println(coverageExample(new(DSRecord).GetSize()))
	// Output: 1
}

func ExampleNSRecord_GetSize() {
	core.Println(coverageExample(new(NSRecord).GetSize()))
	// Output: 1
}

func ExampleGLUE4Record_GetSize() {
	core.Println(coverageExample(new(GLUE4Record).GetSize()))
	// Output: 1
}

func ExampleGLUE6Record_GetSize() {
	core.Println(coverageExample(new(GLUE6Record).GetSize()))
	// Output: 1
}

func ExampleSYNTH4Record_GetSize() {
	core.Println(coverageExample(new(SYNTH4Record).GetSize()))
	// Output: 1
}

func ExampleSYNTH6Record_GetSize() {
	core.Println(coverageExample(new(SYNTH6Record).GetSize()))
	// Output: 1
}

func ExampleTXTRecord_GetSize() {
	core.Println(coverageExample(new(TXTRecord).GetSize()))
	// Output: 1
}

func ExampleDecodeResource() {
	core.Println(coverageExample(DecodeResource([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetDecodeResource() {
	core.Println(coverageExample(GetDecodeResource([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleDSRecord_GetJSON() {
	core.Println(coverageExample(new(DSRecord).GetJSON()))
	// Output: 1
}

func ExampleDSRecord_FromJSON() {
	core.Println(coverageExample(new(DSRecord).FromJSON(coverageValue[DSRecordJSON]("good"))))
	// Output: 1
}

func ExampleNSRecord_GetJSON() {
	core.Println(coverageExample(new(NSRecord).GetJSON()))
	// Output: 1
}

func ExampleNSRecord_FromJSON() {
	core.Println(coverageExample(new(NSRecord).FromJSON(coverageValue[NSRecordJSON]("good"))))
	// Output: 1
}

func ExampleGLUE4Record_GetJSON() {
	core.Println(coverageExample(new(GLUE4Record).GetJSON()))
	// Output: 1
}

func ExampleGLUE4Record_FromJSON() {
	core.Println(coverageExample(new(GLUE4Record).FromJSON(coverageValue[GLUE4RecordJSON]("good"))))
	// Output: 1
}

func ExampleGLUE6Record_GetJSON() {
	core.Println(coverageExample(new(GLUE6Record).GetJSON()))
	// Output: 1
}

func ExampleGLUE6Record_FromJSON() {
	core.Println(coverageExample(new(GLUE6Record).FromJSON(coverageValue[GLUE6RecordJSON]("good"))))
	// Output: 1
}

func ExampleSYNTH4Record_GetJSON() {
	core.Println(coverageExample(new(SYNTH4Record).GetJSON()))
	// Output: 1
}

func ExampleSYNTH4Record_FromJSON() {
	core.Println(coverageExample(new(SYNTH4Record).FromJSON(coverageValue[SYNTH4RecordJSON]("good"))))
	// Output: 1
}

func ExampleSYNTH6Record_GetJSON() {
	core.Println(coverageExample(new(SYNTH6Record).GetJSON()))
	// Output: 1
}

func ExampleSYNTH6Record_FromJSON() {
	core.Println(coverageExample(new(SYNTH6Record).FromJSON(coverageValue[SYNTH6RecordJSON]("good"))))
	// Output: 1
}

func ExampleTXTRecord_GetJSON() {
	core.Println(coverageExample(new(TXTRecord).GetJSON()))
	// Output: 1
}

func ExampleTXTRecord_FromJSON() {
	core.Println(coverageExample(new(TXTRecord).FromJSON(coverageValue[TXTRecordJSON]("good"))))
	// Output: 1
}

// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleRawMessage_MarshalJSON() {
	fn := ax7ValueOf(new(RawMessage).MarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleRawMessage_UnmarshalJSON() {
	fn := ax7ValueOf(new(RawMessage).UnmarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetSize() {
	fn := ax7ValueOf(new(Resource).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewResource() {
	fn := ax7ValueOf(NewResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewResource() {
	fn := ax7ValueOf(GetNewResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDSRecord_Type() {
	fn := ax7ValueOf(new(DSRecord).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDSRecord_GetType() {
	fn := ax7ValueOf(new(DSRecord).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNSRecord_Type() {
	fn := ax7ValueOf(new(NSRecord).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNSRecord_GetType() {
	fn := ax7ValueOf(new(NSRecord).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNSRecord_GetNS() {
	fn := ax7ValueOf(new(NSRecord).GetNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE4Record_Type() {
	fn := ax7ValueOf(new(GLUE4Record).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE4Record_GetType() {
	fn := ax7ValueOf(new(GLUE4Record).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE4Record_GetNS() {
	fn := ax7ValueOf(new(GLUE4Record).GetNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE6Record_Type() {
	fn := ax7ValueOf(new(GLUE6Record).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE6Record_GetType() {
	fn := ax7ValueOf(new(GLUE6Record).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE6Record_GetNS() {
	fn := ax7ValueOf(new(GLUE6Record).GetNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_Type() {
	fn := ax7ValueOf(new(SYNTH4Record).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_GetType() {
	fn := ax7ValueOf(new(SYNTH4Record).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_GetNS() {
	fn := ax7ValueOf(new(SYNTH4Record).GetNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_Type() {
	fn := ax7ValueOf(new(SYNTH6Record).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_GetType() {
	fn := ax7ValueOf(new(SYNTH6Record).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_GetNS() {
	fn := ax7ValueOf(new(SYNTH6Record).GetNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTXTRecord_Type() {
	fn := ax7ValueOf(new(TXTRecord).Type)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTXTRecord_GetType() {
	fn := ax7ValueOf(new(TXTRecord).GetType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_NS() {
	fn := ax7ValueOf(new(SYNTH4Record).NS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_NS() {
	fn := ax7ValueOf(new(SYNTH6Record).NS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_HasType() {
	fn := ax7ValueOf(new(Resource).HasType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetHasType() {
	fn := ax7ValueOf(new(Resource).GetHasType)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_HasNS() {
	fn := ax7ValueOf(new(Resource).HasNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetHasNS() {
	fn := ax7ValueOf(new(Resource).GetHasNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_HasDS() {
	fn := ax7ValueOf(new(Resource).HasDS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetHasDS() {
	fn := ax7ValueOf(new(Resource).GetHasDS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToNS() {
	fn := ax7ValueOf(new(Resource).ToNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToNS() {
	fn := ax7ValueOf(new(Resource).GetToNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToGlue() {
	fn := ax7ValueOf(new(Resource).ToGlue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToGlue() {
	fn := ax7ValueOf(new(Resource).GetToGlue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToDS() {
	fn := ax7ValueOf(new(Resource).ToDS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToDS() {
	fn := ax7ValueOf(new(Resource).GetToDS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToTXT() {
	fn := ax7ValueOf(new(Resource).ToTXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToTXT() {
	fn := ax7ValueOf(new(Resource).GetToTXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToZone() {
	fn := ax7ValueOf(new(Resource).ToZone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToZone() {
	fn := ax7ValueOf(new(Resource).GetToZone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToReferral() {
	fn := ax7ValueOf(new(Resource).ToReferral)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToReferral() {
	fn := ax7ValueOf(new(Resource).GetToReferral)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToDNS() {
	fn := ax7ValueOf(new(Resource).ToDNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToDNS() {
	fn := ax7ValueOf(new(Resource).GetToDNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_ToNSEC() {
	fn := ax7ValueOf(new(Resource).ToNSEC)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetToNSEC() {
	fn := ax7ValueOf(new(Resource).GetToNSEC)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_Encode() {
	fn := ax7ValueOf(new(Resource).Encode)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResource_GetEncode() {
	fn := ax7ValueOf(new(Resource).GetEncode)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResource_Decode() {
	fn := ax7ValueOf(new(Resource).Decode)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetDecode() {
	fn := ax7ValueOf(new(Resource).GetDecode)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_GetJSON() {
	fn := ax7ValueOf(new(Resource).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResource_FromJSON() {
	fn := ax7ValueOf(new(Resource).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDSRecord_GetSize() {
	fn := ax7ValueOf(new(DSRecord).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNSRecord_GetSize() {
	fn := ax7ValueOf(new(NSRecord).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE4Record_GetSize() {
	fn := ax7ValueOf(new(GLUE4Record).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE6Record_GetSize() {
	fn := ax7ValueOf(new(GLUE6Record).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_GetSize() {
	fn := ax7ValueOf(new(SYNTH4Record).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_GetSize() {
	fn := ax7ValueOf(new(SYNTH6Record).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTXTRecord_GetSize() {
	fn := ax7ValueOf(new(TXTRecord).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDecodeResource() {
	fn := ax7ValueOf(DecodeResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetDecodeResource() {
	fn := ax7ValueOf(GetDecodeResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleDSRecord_GetJSON() {
	fn := ax7ValueOf(new(DSRecord).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDSRecord_FromJSON() {
	fn := ax7ValueOf(new(DSRecord).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNSRecord_GetJSON() {
	fn := ax7ValueOf(new(NSRecord).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNSRecord_FromJSON() {
	fn := ax7ValueOf(new(NSRecord).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE4Record_GetJSON() {
	fn := ax7ValueOf(new(GLUE4Record).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE4Record_FromJSON() {
	fn := ax7ValueOf(new(GLUE4Record).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE6Record_GetJSON() {
	fn := ax7ValueOf(new(GLUE6Record).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGLUE6Record_FromJSON() {
	fn := ax7ValueOf(new(GLUE6Record).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_GetJSON() {
	fn := ax7ValueOf(new(SYNTH4Record).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH4Record_FromJSON() {
	fn := ax7ValueOf(new(SYNTH4Record).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_GetJSON() {
	fn := ax7ValueOf(new(SYNTH6Record).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleSYNTH6Record_FromJSON() {
	fn := ax7ValueOf(new(SYNTH6Record).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTXTRecord_GetJSON() {
	fn := ax7ValueOf(new(TXTRecord).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTXTRecord_FromJSON() {
	fn := ax7ValueOf(new(TXTRecord).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

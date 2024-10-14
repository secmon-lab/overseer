// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/m-mizutani/opac"
	"github.com/open-policy-agent/opa/ast"
	"github.com/secmon-as-code/overseer/pkg/interfaces"
	"io"
	"sync"
)

// Ensure, that CloudStorageClientMock does implement interfaces.CloudStorageClient.
// If this is not the case, regenerate this file with moq.
var _ interfaces.CloudStorageClient = &CloudStorageClientMock{}

// CloudStorageClientMock is a mock implementation of interfaces.CloudStorageClient.
//
//	func TestSomethingThatUsesCloudStorageClient(t *testing.T) {
//
//		// make and configure a mocked interfaces.CloudStorageClient
//		mockedCloudStorageClient := &CloudStorageClientMock{
//			GetObjectFunc: func(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
//				panic("mock out the GetObject method")
//			},
//			PutObjectFunc: func(ctx context.Context, bucketName string, objectName string) (io.WriteCloser, error) {
//				panic("mock out the PutObject method")
//			},
//		}
//
//		// use mockedCloudStorageClient in code that requires interfaces.CloudStorageClient
//		// and then make assertions.
//
//	}
type CloudStorageClientMock struct {
	// GetObjectFunc mocks the GetObject method.
	GetObjectFunc func(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error)

	// PutObjectFunc mocks the PutObject method.
	PutObjectFunc func(ctx context.Context, bucketName string, objectName string) (io.WriteCloser, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetObject holds details about calls to the GetObject method.
		GetObject []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// BucketName is the bucketName argument value.
			BucketName string
			// ObjectName is the objectName argument value.
			ObjectName string
		}
		// PutObject holds details about calls to the PutObject method.
		PutObject []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// BucketName is the bucketName argument value.
			BucketName string
			// ObjectName is the objectName argument value.
			ObjectName string
		}
	}
	lockGetObject sync.RWMutex
	lockPutObject sync.RWMutex
}

// GetObject calls GetObjectFunc.
func (mock *CloudStorageClientMock) GetObject(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
	if mock.GetObjectFunc == nil {
		panic("CloudStorageClientMock.GetObjectFunc: method is nil but CloudStorageClient.GetObject was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		BucketName string
		ObjectName string
	}{
		Ctx:        ctx,
		BucketName: bucketName,
		ObjectName: objectName,
	}
	mock.lockGetObject.Lock()
	mock.calls.GetObject = append(mock.calls.GetObject, callInfo)
	mock.lockGetObject.Unlock()
	return mock.GetObjectFunc(ctx, bucketName, objectName)
}

// GetObjectCalls gets all the calls that were made to GetObject.
// Check the length with:
//
//	len(mockedCloudStorageClient.GetObjectCalls())
func (mock *CloudStorageClientMock) GetObjectCalls() []struct {
	Ctx        context.Context
	BucketName string
	ObjectName string
} {
	var calls []struct {
		Ctx        context.Context
		BucketName string
		ObjectName string
	}
	mock.lockGetObject.RLock()
	calls = mock.calls.GetObject
	mock.lockGetObject.RUnlock()
	return calls
}

// PutObject calls PutObjectFunc.
func (mock *CloudStorageClientMock) PutObject(ctx context.Context, bucketName string, objectName string) (io.WriteCloser, error) {
	if mock.PutObjectFunc == nil {
		panic("CloudStorageClientMock.PutObjectFunc: method is nil but CloudStorageClient.PutObject was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		BucketName string
		ObjectName string
	}{
		Ctx:        ctx,
		BucketName: bucketName,
		ObjectName: objectName,
	}
	mock.lockPutObject.Lock()
	mock.calls.PutObject = append(mock.calls.PutObject, callInfo)
	mock.lockPutObject.Unlock()
	return mock.PutObjectFunc(ctx, bucketName, objectName)
}

// PutObjectCalls gets all the calls that were made to PutObject.
// Check the length with:
//
//	len(mockedCloudStorageClient.PutObjectCalls())
func (mock *CloudStorageClientMock) PutObjectCalls() []struct {
	Ctx        context.Context
	BucketName string
	ObjectName string
} {
	var calls []struct {
		Ctx        context.Context
		BucketName string
		ObjectName string
	}
	mock.lockPutObject.RLock()
	calls = mock.calls.PutObject
	mock.lockPutObject.RUnlock()
	return calls
}

// Ensure, that BigQueryClientMock does implement interfaces.BigQueryClient.
// If this is not the case, regenerate this file with moq.
var _ interfaces.BigQueryClient = &BigQueryClientMock{}

// BigQueryClientMock is a mock implementation of interfaces.BigQueryClient.
//
//	func TestSomethingThatUsesBigQueryClient(t *testing.T) {
//
//		// make and configure a mocked interfaces.BigQueryClient
//		mockedBigQueryClient := &BigQueryClientMock{
//			QueryFunc: func(ctx context.Context, query string) (interfaces.BigQueryIterator, error) {
//				panic("mock out the Query method")
//			},
//		}
//
//		// use mockedBigQueryClient in code that requires interfaces.BigQueryClient
//		// and then make assertions.
//
//	}
type BigQueryClientMock struct {
	// QueryFunc mocks the Query method.
	QueryFunc func(ctx context.Context, query string) (interfaces.BigQueryIterator, error)

	// calls tracks calls to the methods.
	calls struct {
		// Query holds details about calls to the Query method.
		Query []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Query is the query argument value.
			Query string
		}
	}
	lockQuery sync.RWMutex
}

// Query calls QueryFunc.
func (mock *BigQueryClientMock) Query(ctx context.Context, query string) (interfaces.BigQueryIterator, error) {
	if mock.QueryFunc == nil {
		panic("BigQueryClientMock.QueryFunc: method is nil but BigQueryClient.Query was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Query string
	}{
		Ctx:   ctx,
		Query: query,
	}
	mock.lockQuery.Lock()
	mock.calls.Query = append(mock.calls.Query, callInfo)
	mock.lockQuery.Unlock()
	return mock.QueryFunc(ctx, query)
}

// QueryCalls gets all the calls that were made to Query.
// Check the length with:
//
//	len(mockedBigQueryClient.QueryCalls())
func (mock *BigQueryClientMock) QueryCalls() []struct {
	Ctx   context.Context
	Query string
} {
	var calls []struct {
		Ctx   context.Context
		Query string
	}
	mock.lockQuery.RLock()
	calls = mock.calls.Query
	mock.lockQuery.RUnlock()
	return calls
}

// Ensure, that BigQueryIteratorMock does implement interfaces.BigQueryIterator.
// If this is not the case, regenerate this file with moq.
var _ interfaces.BigQueryIterator = &BigQueryIteratorMock{}

// BigQueryIteratorMock is a mock implementation of interfaces.BigQueryIterator.
//
//	func TestSomethingThatUsesBigQueryIterator(t *testing.T) {
//
//		// make and configure a mocked interfaces.BigQueryIterator
//		mockedBigQueryIterator := &BigQueryIteratorMock{
//			NextFunc: func(dst interface{}) error {
//				panic("mock out the Next method")
//			},
//		}
//
//		// use mockedBigQueryIterator in code that requires interfaces.BigQueryIterator
//		// and then make assertions.
//
//	}
type BigQueryIteratorMock struct {
	// NextFunc mocks the Next method.
	NextFunc func(dst interface{}) error

	// calls tracks calls to the methods.
	calls struct {
		// Next holds details about calls to the Next method.
		Next []struct {
			// Dst is the dst argument value.
			Dst interface{}
		}
	}
	lockNext sync.RWMutex
}

// Next calls NextFunc.
func (mock *BigQueryIteratorMock) Next(dst interface{}) error {
	if mock.NextFunc == nil {
		panic("BigQueryIteratorMock.NextFunc: method is nil but BigQueryIterator.Next was just called")
	}
	callInfo := struct {
		Dst interface{}
	}{
		Dst: dst,
	}
	mock.lockNext.Lock()
	mock.calls.Next = append(mock.calls.Next, callInfo)
	mock.lockNext.Unlock()
	return mock.NextFunc(dst)
}

// NextCalls gets all the calls that were made to Next.
// Check the length with:
//
//	len(mockedBigQueryIterator.NextCalls())
func (mock *BigQueryIteratorMock) NextCalls() []struct {
	Dst interface{}
} {
	var calls []struct {
		Dst interface{}
	}
	mock.lockNext.RLock()
	calls = mock.calls.Next
	mock.lockNext.RUnlock()
	return calls
}

// Ensure, that PubSubClientMock does implement interfaces.PubSubClient.
// If this is not the case, regenerate this file with moq.
var _ interfaces.PubSubClient = &PubSubClientMock{}

// PubSubClientMock is a mock implementation of interfaces.PubSubClient.
//
//	func TestSomethingThatUsesPubSubClient(t *testing.T) {
//
//		// make and configure a mocked interfaces.PubSubClient
//		mockedPubSubClient := &PubSubClientMock{
//			PublishFunc: func(ctx context.Context, topic string, data []byte) error {
//				panic("mock out the Publish method")
//			},
//		}
//
//		// use mockedPubSubClient in code that requires interfaces.PubSubClient
//		// and then make assertions.
//
//	}
type PubSubClientMock struct {
	// PublishFunc mocks the Publish method.
	PublishFunc func(ctx context.Context, topic string, data []byte) error

	// calls tracks calls to the methods.
	calls struct {
		// Publish holds details about calls to the Publish method.
		Publish []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Topic is the topic argument value.
			Topic string
			// Data is the data argument value.
			Data []byte
		}
	}
	lockPublish sync.RWMutex
}

// Publish calls PublishFunc.
func (mock *PubSubClientMock) Publish(ctx context.Context, topic string, data []byte) error {
	if mock.PublishFunc == nil {
		panic("PubSubClientMock.PublishFunc: method is nil but PubSubClient.Publish was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Topic string
		Data  []byte
	}{
		Ctx:   ctx,
		Topic: topic,
		Data:  data,
	}
	mock.lockPublish.Lock()
	mock.calls.Publish = append(mock.calls.Publish, callInfo)
	mock.lockPublish.Unlock()
	return mock.PublishFunc(ctx, topic, data)
}

// PublishCalls gets all the calls that were made to Publish.
// Check the length with:
//
//	len(mockedPubSubClient.PublishCalls())
func (mock *PubSubClientMock) PublishCalls() []struct {
	Ctx   context.Context
	Topic string
	Data  []byte
} {
	var calls []struct {
		Ctx   context.Context
		Topic string
		Data  []byte
	}
	mock.lockPublish.RLock()
	calls = mock.calls.Publish
	mock.lockPublish.RUnlock()
	return calls
}

// Ensure, that PolicyClientMock does implement interfaces.PolicyClient.
// If this is not the case, regenerate this file with moq.
var _ interfaces.PolicyClient = &PolicyClientMock{}

// PolicyClientMock is a mock implementation of interfaces.PolicyClient.
//
//	func TestSomethingThatUsesPolicyClient(t *testing.T) {
//
//		// make and configure a mocked interfaces.PolicyClient
//		mockedPolicyClient := &PolicyClientMock{
//			AnnotationSetFunc: func() *ast.AnnotationSet {
//				panic("mock out the AnnotationSet method")
//			},
//			QueryFunc: func(ctx context.Context, query string, input any, output any, options ...opac.QueryOption) error {
//				panic("mock out the Query method")
//			},
//		}
//
//		// use mockedPolicyClient in code that requires interfaces.PolicyClient
//		// and then make assertions.
//
//	}
type PolicyClientMock struct {
	// AnnotationSetFunc mocks the AnnotationSet method.
	AnnotationSetFunc func() *ast.AnnotationSet

	// QueryFunc mocks the Query method.
	QueryFunc func(ctx context.Context, query string, input any, output any, options ...opac.QueryOption) error

	// calls tracks calls to the methods.
	calls struct {
		// AnnotationSet holds details about calls to the AnnotationSet method.
		AnnotationSet []struct {
		}
		// Query holds details about calls to the Query method.
		Query []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Query is the query argument value.
			Query string
			// Input is the input argument value.
			Input any
			// Output is the output argument value.
			Output any
			// Options is the options argument value.
			Options []opac.QueryOption
		}
	}
	lockAnnotationSet sync.RWMutex
	lockQuery         sync.RWMutex
}

// AnnotationSet calls AnnotationSetFunc.
func (mock *PolicyClientMock) AnnotationSet() *ast.AnnotationSet {
	if mock.AnnotationSetFunc == nil {
		panic("PolicyClientMock.AnnotationSetFunc: method is nil but PolicyClient.AnnotationSet was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAnnotationSet.Lock()
	mock.calls.AnnotationSet = append(mock.calls.AnnotationSet, callInfo)
	mock.lockAnnotationSet.Unlock()
	return mock.AnnotationSetFunc()
}

// AnnotationSetCalls gets all the calls that were made to AnnotationSet.
// Check the length with:
//
//	len(mockedPolicyClient.AnnotationSetCalls())
func (mock *PolicyClientMock) AnnotationSetCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAnnotationSet.RLock()
	calls = mock.calls.AnnotationSet
	mock.lockAnnotationSet.RUnlock()
	return calls
}

// Query calls QueryFunc.
func (mock *PolicyClientMock) Query(ctx context.Context, query string, input any, output any, options ...opac.QueryOption) error {
	if mock.QueryFunc == nil {
		panic("PolicyClientMock.QueryFunc: method is nil but PolicyClient.Query was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Query   string
		Input   any
		Output  any
		Options []opac.QueryOption
	}{
		Ctx:     ctx,
		Query:   query,
		Input:   input,
		Output:  output,
		Options: options,
	}
	mock.lockQuery.Lock()
	mock.calls.Query = append(mock.calls.Query, callInfo)
	mock.lockQuery.Unlock()
	return mock.QueryFunc(ctx, query, input, output, options...)
}

// QueryCalls gets all the calls that were made to Query.
// Check the length with:
//
//	len(mockedPolicyClient.QueryCalls())
func (mock *PolicyClientMock) QueryCalls() []struct {
	Ctx     context.Context
	Query   string
	Input   any
	Output  any
	Options []opac.QueryOption
} {
	var calls []struct {
		Ctx     context.Context
		Query   string
		Input   any
		Output  any
		Options []opac.QueryOption
	}
	mock.lockQuery.RLock()
	calls = mock.calls.Query
	mock.lockQuery.RUnlock()
	return calls
}

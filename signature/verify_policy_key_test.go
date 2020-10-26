/*
 * Copyright 2020 Guardtime, Inc.
 *
 * This file is part of the Guardtime client SDK.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES, CONDITIONS, OR OTHER LICENSES OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 * "Guardtime" and "KSI" are trademarks or registered trademarks of
 * Guardtime, Inc., and no license to trademarks is granted; Guardtime
 * reserves and retains all trademark rights.
 */

package signature

import (
	"path/filepath"
	"testing"

	"github.com/vincentracine/goksi/log"
	"github.com/vincentracine/goksi/publications"
	"github.com/vincentracine/goksi/signature/verify/result"
	"github.com/vincentracine/goksi/test"
)

func TestUnitKeyBasedVerificationPolicy(t *testing.T) {

	logger, defFunc, err := test.InitLogger(t, testLogDir, log.DEBUG, t.Name())
	if err != nil {
		t.Fatal("Failed to initialize logger: ", err)
	}
	defer defFunc()
	// Apply logger.
	log.SetLogger(logger)

	test.Suite{
		{Func: testKeyPolicyVerifySignature},
	}.Runner(t)
}

func testKeyPolicyVerifySignature(t *testing.T, _ ...interface{}) {

	var (
		testSigFile = filepath.Join(testSigDir, "ok-sig-2014-08-01.1.ksig")
		testPubFile = filepath.Join(testPubDir, "publications.bin")
		testCrtFile = filepath.Join(testCrtDir, "mock.crt")
		testPolicy  = KeyBasedVerificationPolicy
	)

	sig, err := New(BuildNoVerify(BuildFromFile(testSigFile)))
	if err != nil {
		t.Fatal("Failed to create ksi signature: ", err)
	}

	pubFile, err := publications.NewFile(publications.FileFromFile(testPubFile))
	if err != nil {
		t.Fatal("Failed to initialize publications file: ", err)
	}

	pfh, err := publications.NewFileHandler(
		publications.FileHandlerSetFile(pubFile),
		publications.FileHandlerSetTrustedCertificateFromFilePem(testCrtFile),
		publications.FileHandlerSetFileCertConstraint(publications.OidEmail, "publications@guardtime.com"),
	)
	if err != nil {
		t.Fatal("Failed to initialize publications file handler: ", err)
	}

	verCtx, err := NewVerificationContext(sig,
		VerCtxOptPublicationsFileHandler(pfh),
	)
	if err != nil {
		t.Fatal("Failed to create verification context: ", err)
	}

	res, err := testPolicy.Verify(verCtx)
	if err != nil {
		t.Fatal(testPolicy, " returned error: ", err)
	}
	if res != result.OK {
		verRes, _ := verCtx.Result()
		t.Error(testPolicy, "Verify result: ", res)
		t.Error(testPolicy, "Final  result: ", verRes.FinalResult())
	}
}

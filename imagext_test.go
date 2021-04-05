/*
 *          Copyright 2021, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *      (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package imagext

import (
	"testing"
)

func TestGray(t *testing.T) {
	// 0,2126 * R + 0,7152 * B + 0,0722 * G
	gray := Gray(0, 0, 0)
	if gray != 0 {
		t.Error(gray)
	}
	gray = Gray(255, 255, 255)
	if gray != 255 {
		t.Error(gray)
	}
	gray = Gray(128, 128, 128)
	if gray != 128 {
		t.Error(gray)
	}
	gray = Gray(64, 128, 64)
	if gray != 109 {
		t.Error(gray)
	}
}

func TestCMYKToGray(t *testing.T) {
	// 0.299 * R + 0.587 * G + 0.114 * B
	gray := cmykToGray(0, 0, 0, 0)
	if gray != 255 {
		t.Error(gray)
	}
	gray = cmykToGray(255, 255, 255, 0)
	if gray != 0 {
		t.Error(gray)
	}
	gray = cmykToGray(0, 0, 0, 255)
	if gray != 0 {
		t.Error(gray)
	}
	gray = cmykToGray(127, 127, 127, 127)
	if gray != 64 {
		t.Error(gray)
	}
	gray = cmykToGray(0, 0, 0, 127)
	if gray != 128 {
		t.Error(gray)
	}
	gray = cmykToGray(0, 0, 0, 128)
	if gray != 127 {
		t.Error(gray)
	}
}

func TestMedian(t *testing.T) {
	hist := make([]uint8, 10)
	histZero := make([]uint8, 10)
	hist[0], hist[3], hist[5], hist[6] = 2, 7, 1, 3
	med := median(hist)
	if med != 3 {
		t.Error(med)
	}
	copy(hist, histZero)
	hist[0], hist[3], hist[5], hist[6] = 7, 2, 1, 3
	med = median(hist)
	if med != 0 {
		t.Error(med)
	}
	copy(hist, histZero)
	hist[0], hist[3], hist[5], hist[6], hist[9] = 7, 2, 1, 3, 10
	med = median(hist)
	if med != 6 {
		t.Error(med)
	}
}

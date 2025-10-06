/*
 *    Copyright (c) 2024 wslyyy
 *
 *    Permission is hereby granted, free of charge, to any person obtaining a copy
 *    of this software and associated documentation files (the "Software"), to deal
 *    in the Software without restriction, including without limitation the rights
 *    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *    copies of the Software, and to permit persons to whom the Software is
 *    furnished to do so, subject to the following conditions:
 *
 *    The above copyright notice and this permission notice shall be included in all
 *    copies or substantial portions of the Software.
 *
 *    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *    SOFTWARE.
 */

package innertube

func Filter(dictionary map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range dictionary {
		if value != nil {
			switch v := value.(type) {
			case *int:
				if v != nil {
					result[key] = value
				}
			case *string:
				if v != nil {
					result[key] = value
				}
			default:
				result[key] = value
			}
		}
	}
	return result
}

func Contextualise(clientContext ClientContext, data map[string]interface{}) map[string]interface{} {
	if _, ok := data["context"]; !ok {
		data["context"] = make(map[string]interface{})
	}
	if _, ok := data["context"].(map[string]interface{})["client"]; !ok {
		data["context"].(map[string]interface{})["client"] = make(map[string]interface{})
	}

	clientData := clientContext.Context()
	for key, value := range clientData {
		data["context"].(map[string]interface{})["client"].(map[string]interface{})[key] = value
	}

	return data
}

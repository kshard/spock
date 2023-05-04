/*

  Knowledge Graph: SPOCK
  Copyright (C) 2016 - 2023 Dmitry Kolesnikov

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published
  by the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package ephemeral

import (
	"fmt"

	"github.com/kshard/spock"
)

type notSupported struct{ spock.Pattern }

func (err notSupported) Error() string { return fmt.Sprintf("not supported %s", err.Pattern.Dump()) }
func (notSupported) NotSupported()     {}

func (store *Store) streamSPO(q spock.Pattern) (spock.Stream, error) {
	return newIterator[s, p, o](querySPO(q), store.spo), nil
}

func (store *Store) streamSOP(q spock.Pattern) (spock.Stream, error) {
	return newIterator[s, o, p](querySOP(q), store.sop), nil
}

func (store *Store) streamPSO(q spock.Pattern) (spock.Stream, error) {
	return newIterator[p, s, o](queryPSO(q), store.pso), nil
}

func (store *Store) streamPOS(q spock.Pattern) (spock.Stream, error) {
	return newIterator[p, o, s](queryPOS(q), store.pos), nil
}

func (store *Store) streamOSP(q spock.Pattern) (spock.Stream, error) {
	return newIterator[o, s, p](queryOSP(q), store.osp), nil
}

func (store *Store) streamOPS(q spock.Pattern) (spock.Stream, error) {
	return newIterator[o, p, s](queryOPS(q), store.ops), nil
}

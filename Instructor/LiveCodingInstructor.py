import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

LCI_DEQUE_PATH = "deque"
LCI_BROWNIE_PATH = "brownie"

LCI_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "LiveCodingInstructorInfo")

def lci_set_attr(attr):
	def foo(value):
		try:
			with open(LCI_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
		json_obj[attr] = value
		with open(LCI_FILE, 'w') as f:
			f.write(json.dumps(json_obj))

	return foo


def lci_get_attr():
	try:
		with open(LCI_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set server address and passcode.")
		return None
	if 'Address' not in json_obj:
		sublime.message_dialog("Please set server address.")
		return None
	if 'Passcode' not in json_obj:
		sublime.message_dialog("Please set passcode.")
		return None
	if not json_obj['Address'].startswith('http'):
		json_obj['Address'] = 'http://' + json_obj['Address']
	return json_obj


class LciSetServerAddressCommand(sublime_plugin.WindowCommand):
	def run(self):
		info = lci_get_attr()
		sublime.active_window().show_input_panel('Server address: ', info['Address'], lci_set_attr('Address'), None, None)

class LciSetPasscodeCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Passcode: ', '', lci_set_attr('Passcode'), None, None)


class LciGetCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		global CUR_USER
		info = lci_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Address'], LCI_DEQUE_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		req = urllib.request.Request(url, data)
		try:
			with urllib.request.urlopen(req) as response:
				json_obj = json.loads(response.read().decode(encoding="utf-8"))
				self.view.replace(edit, sublime.Region(0, self.view.size()), json_obj['Body'])
				sublime.status_message(json_obj['User'] + ", " + str(json_obj['N']) + " entries left")
		except urllib.error.HTTPError:
			sublime.message_dialog("HTTP error: possibly due to incorrect passcode.")
		except urllib.error.URLError:
			sublime.message_dialog("URL error: reset server address.")


class LciAwardPointCommand(sublime_plugin.WindowCommand):
	def run(self):
		info = lci_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Address'], LCI_BROWNIE_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		req = urllib.request.Request(url, data)
		try:
			with urllib.request.urlopen(req) as response:
				user = response.read().decode(encoding="utf-8")
				if user:
					sublime.message_dialog("+1 brownie for " + user)
				else:
					sublime.message_dialog("No current user.")
		except urllib.error.HTTPError:
			sublime.message_dialog("HTTP error: possibly due to incorrect passcode.")
		except urllib.error.URLError:
			sublime.message_dialog("URL error: reset server address.")



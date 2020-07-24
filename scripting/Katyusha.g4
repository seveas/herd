grammar Katyusha;

// Keep this  the top, as it gobbles up everything after it
RUN: 'run' ~('\n')* ;
SB_OPEN: '[' ;
CB_OPEN: '{' ;
SET: 'set' ;
ADD: 'add' ;
REMOVE: 'remove' ;
LIST: 'list' ;
HOSTS: 'hosts' ;
DURATION: ( '-'? [0-9]+ ( '.' [0-9]+ )? [smh] )+ ;
NUMBER: '0x'?[0-9]+ ;
IDENTIFIER: ( [a-zA-Z_][-a-zA-Z_.:0-9]*[a-zA-Z_0-9] | [a-zA-Z] );
GLOB: [-a-zA-Z.0-9*?]+ ;
EQUALS: '==' ;
MATCHES: '=~' ;
NOT_EQUALS: '!=';
NOT_MATCHES: '!~';
STRING
 : '"' ( '\\' . | ~[\\\r\n\f"] )* '"'
 ;
REGEXP
 : '/' ( '\\' . | ~[\\\r\n\f/] )* '/'
 ;

fragment COMMENT: '#' ~('\n')+;
fragment SPACES: ' '+ ;

SKIP_ : ( SPACES | COMMENT ) -> skip ;

prog : line* EOF ;
line : ( run | set | add | remove | list )? '\n' ;
run : RUN ;
set: SET (varname=IDENTIFIER varvalue=scalar)? ;
add: ADD HOSTS ( glob=(GLOB|IDENTIFIER) filters=filter* | filters=filter+ );
remove: REMOVE HOSTS ( glob=(GLOB|IDENTIFIER) filters=filter* | filters=filter+ );
list: LIST HOSTS opts=hash? ;
filter: key=IDENTIFIER ( comp=( EQUALS | NOT_EQUALS ) val=scalar | comp=( MATCHES | NOT_MATCHES ) rx=REGEXP );
scalar: NUMBER | STRING | DURATION | IDENTIFIER ;
value: scalar | array | hash ;
array: ( '[' ']' | '[' value (',' value)* ']' );
hash: ( '{' '}' | '{' IDENTIFIER ':' value (',' IDENTIFIER ':' value)* '}' );
